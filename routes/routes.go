package routes

import (
	"context"
	"net/http"
	"purchaseOrderSystem/auth"
	"purchaseOrderSystem/components"
	"purchaseOrderSystem/services"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This router is for all urls that are publically accessible
func PublicRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	userService := services.NewUserService(db)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFound)
		components.Hello("Aman").Render(context.Background(), w)
	})
	// FIX: not the best way to serve static pages. I have to find a better way
	signupPage := template.Must(template.ParseFiles("public/signup.html"))
	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAlreadyReported)
		signupPage.Execute(w, nil)
	})
	r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := userService.Signup(email, username, password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// FIX: Need to make sure that the user is correctly routed to the homepage of their respective role
		return
	})

	return r
}

// UserRouter is for routes for all Users
func UserRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.UserOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello User"))
	})
	return r
}

// ManagerRouter is for routes for all managers
func ManagerRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.ManagerOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Manager"))
	})
	return r
}

// adminRouter is for routes for all admins
func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello admin"))
	})
	return r
}
