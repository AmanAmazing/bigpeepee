package routes

import (
	"context"
	"net/http"
	"purchaseOrderSystem/auth"
	"purchaseOrderSystem/components"
	"purchaseOrderSystem/services"
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
	homePage := template.Must(template.ParseFiles("public/home.html"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		components.Hello("Aman").Render(context.Background(), w)
	})

	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

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
	formPage := template.Must(template.ParseFiles("components/poform.html"))
	userService := services.NewUserService(db)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello User"))
	})
	r.Get("/form", func(w http.ResponseWriter, r *http.Request) {
		formPage.Execute(w, nil)
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
	r.Post("/form/submit", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// Handle the error
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// // Get the priority value
		// priority := r.PostForm.Get("priority")
		//
		// // Get the item data
		// names := r.PostForm["name[]"]
		// suppliers := r.PostForm["supplier[]"]
		// nominals := r.PostForm["nominal[]"]
		// products := r.PostForm["product[]"]
		// unitPrices := r.PostForm["unit_price[]"]
		// quantities := r.PostForm["quantity[]"]
		// links := r.PostForm["link[]"]
		//
		// // Process the form data
		// for i := 0; i < len(names); i++ {
		// 	name := names[i]
		// 	supplier := suppliers[i]
		// 	nominal := nominals[i]
		// 	product := products[i]
		// 	unitPrice := unitPrices[i]
		// 	quantity := quantities[i]
		// 	link := links[i]
		//
		// 	// Perform further processing or store the data in the database
		// 	// ...
		// }
		//
		// Send a response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Form submitted successfully"))
	})

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
