package routes

import (
	"context"
	"fmt"
	"net/http"
	"purchaseOrderSystem/auth"
	"purchaseOrderSystem/components"
	"purchaseOrderSystem/models"
	"purchaseOrderSystem/services"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This router is for all urls that are publically accessible
func PublicRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	userService := services.NewUserService(db)
	// FIX: not the best way to serve static pages. I have to find a better way
	loginPage := template.Must(template.ParseFiles("public/login.html"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		components.Hello("Aman").Render(context.Background(), w)
	})

	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		components.Hello("Something").Render(context.Background(), w)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		loginPage.Execute(w, nil)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// getting the jwt token
		jwtToken, user_role, err := userService.Login(username, password)
		if err != nil {
			if err.Error() == "Invalid username or password" {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    jwtToken,
			HttpOnly: true,
			// secure:   true, // set to true if using https
			SameSite: http.SameSiteStrictMode,
		})

		// redirecting
		switch user_role {
		case "admin":
			http.Redirect(w, r, "/admin/home", http.StatusFound)
		case "manager":
			http.Redirect(w, r, "/user/home", http.StatusFound)
		case "user":
			http.Redirect(w, r, "/user/home", http.StatusFound)
		default:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

	return r
}

// TODO: Need to make the routers check for jwt expiry
// Need a way to refresh the token so the user does not need to sign in again if they have signed in the previous 24 hours

// UserRouter is for routes for all Users
func UserRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(auth.TokenAuth))
	r.Use(auth.AuthMiddleware)
	r.Use(auth.UserMiddlerware())
	// FIX: not the best way to serve static pages. I have to find a better way
	// formPage := template.Must(template.ParseFiles("components/poform.html"))
	userService := services.NewUserService(db)
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		components.UserHome().Render(context.Background(), w)
	})
	r.Get("/form", func(w http.ResponseWriter, r *http.Request) {
		// formPage.Execute(w, nil)
		components.Form().Render(context.Background(), w)
	})

	r.Get("/form/suppliers", func(w http.ResponseWriter, r *http.Request) {
		suppliers, err := userService.GetSuppliers()
		if err != nil {
			println("error occurred fetching suppliers")
		}
		tmpl := template.Must(template.New("suppliers").Parse(`
			{{range .}}
				<option value="{{.ID}}">{{.Name}}</option>
			{{end}}
		`))
		tmpl.Execute(w, suppliers)
	})
	r.Get("/form/nominals", func(w http.ResponseWriter, r *http.Request) {
		nominals, err := userService.GetNominals()
		if err != nil {
			println("error occurred fetching nominals")
		}
		tmpl := template.Must(template.New("nominals").Parse(`
			{{range .}}
				<option value="{{.ID}}">{{.Name}}</option>
			{{end}}
		`))
		tmpl.Execute(w, nominals)
	})
	r.Get("/form/products", func(w http.ResponseWriter, r *http.Request) {
		products, err := userService.GetProducts()
		if err != nil {
			println("error occurred fetching products")
		}
		tmpl := template.Must(template.New("products").Parse(`
			{{range .}}
				<option value="{{.ID}}">{{.Name}}</option>
			{{end}}
		`))
		tmpl.Execute(w, products)
	})
	r.Get("/form/item", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Need to verify the item_count
		item_count_string := r.URL.Query().Get("item_count")
		item_count, err := strconv.Atoi(item_count_string)
		item_count += 1
		if err != nil {
			// TODO: Need to return error
			return
		}
		w.Header().Add("HX-Trigger", "updateItemCountEvet")
		w.Header().Add("HX-Trigger-After-Swap", "updateItemCountEvent")
		components.FormItem(item_count).Render(context.Background(), w)
	})
	r.Post("/form/submit", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// Handle the error
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		item_count, err := strconv.Atoi(r.FormValue("item_count"))
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		priority := r.FormValue("priority")

		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "Failed to gather ID and department ID in the /user/form/submit post path", http.StatusUnauthorized)
			return
		}
		// FIX: need to verify department and userId is there
		user_id := int(claims["userId"].(float64)) // FIX: DUMB
		department := claims["department"].(string)
		role := claims["role"].(string)
		err = userService.SubmitPurchaseOrder(user_id, department, priority, role, item_count, r.Form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Form submitted successfully"))
	})

	r.Get("/repo", func(w http.ResponseWriter, r *http.Request) {
		// need to get all the Purchase order forms raised by someone
		_, claims, err := jwtauth.FromContext(r.Context())
		user_id := int(claims["userId"].(float64)) // FIX: DUMB
		// userId, err := auth.ExtractUserID(r)
		if err != nil {
			// FIX: Need a better server error to throw to the user
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		purchaseOrders, err := userService.GetPurchaseOrdersByUserID(user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		components.Repo(purchaseOrders).Render(context.Background(), w)
		return
	})

	r.Get("/po/{id}", func(w http.ResponseWriter, r *http.Request) {
		requestedPoId := chi.URLParam(r, "id")
		if requestedPoId == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error occurred"))
			return
		}
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			// FIX: Need a better server error to throw to the user
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user_id := int(claims["userId"].(float64)) // FIX: DUMB
		po, err := userService.GetPurchaseOrderById(user_id, requestedPoId)
		if err != nil {
			// TODO: return an error on the repo page itself that the user is not able to access that page to edit.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// render the assoicated page.
		components.PurchaseOrder(po).Render(context.Background(), w)
	})

	r.Get("/po/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		requestedPoId := chi.URLParam(r, "id")
		if requestedPoId == "" {
			// FIX: Need better error responses
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error occurred"))
			return
		}
		// _, claims, err := jwtauth.FromContext(r.Context())
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// user_id = int(claims["user_id"].(float64))
		po, err := userService.GetPurchaseOrderByIdWithoutItems(requestedPoId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		components.EditPurchaseOrder(po).Render(context.Background(), w)

	})

	// posts the edited form. Need to edit the business logic for this.
	r.Put("/po/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		requestedPoId := chi.URLParam(r, "id")
		if requestedPoId == "" {
			// FIX: Need better error responses. Only testing this for now.
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to provide id for po. check url param"))
			return
		}
		// fetching the posted data
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to parse form values"))
			return
		}
		var data models.PurchaseOrder
		data.Title = r.Form.Get("title")
		data.Description = r.Form.Get("description")
		data.Priority = r.Form.Get("priority")
		data.ID, err = strconv.Atoi(requestedPoId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to convert po id from string to integer"))
			return
		}

		if data.IsEmpty() {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Empty form data"))
			return
		}

		_, err = userService.PutPurchaseOrder(data)
		if err != nil {
			// TODO: Redirecting for now but need to send back errors next time
			http.Redirect(w, r, fmt.Sprintf("/user/po/%v", data.ID), http.StatusBadRequest)
			return
		}
		w.Header().Add("HX-Redirect", fmt.Sprintf("/user/po/%v", requestedPoId))
		return
	})

	return r
}

// adminRouter is for routes for all admins
func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(auth.TokenAuth))
	r.Use(auth.AuthMiddleware)
	r.Use(auth.RoleMiddleware("admin"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello admin"))
	})
	return r
}
