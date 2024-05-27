package routes

import (
	"context"
	"net/http"
	"purchaseOrderSystem/auth"
	"purchaseOrderSystem/components"
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
			http.Redirect(w, r, "/admin", http.StatusFound)
		case "manager":
			http.Redirect(w, r, "/manager", http.StatusFound)
		case "user":
			http.Redirect(w, r, "/user", http.StatusFound)
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
	r.Use(auth.RoleMiddleware("user"))
	// FIX: not the best way to serve static pages. I have to find a better way
	// formPage := template.Must(template.ParseFiles("components/poform.html"))
	userService := services.NewUserService(db)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello User"))
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
		user_id := int(claims["userId"].(float64))
		department := claims["department"].(string)
		err = userService.SubmitPurchaseOrder(user_id, department, priority, item_count, r.Form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Form submitted successfully"))
	})

	r.Get("/repo", func(w http.ResponseWriter, r *http.Request) {
		components.Repo().Render(context.Background(), w)
		return
	})
	// r.Get("/po/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	requestedPoId := chi.URLParam(r, "id")
	// 	po, err := userService.GetPurchaseOrderId(requestedPoId)
	// })

	return r
}

// ManagerRouter is for routes for all managers
func ManagerRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(auth.TokenAuth))
	r.Use(auth.AuthMiddleware)
	r.Use(auth.RoleMiddleware("manager"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Manager"))
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
