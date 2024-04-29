package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"purchaseOrderSystem/components"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

// TODO: Add authentcation middleware
// FIX: Need to fix templ support for neovim

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
	fmt.Println("Started server on port", os.Getenv("PORT"))

	r := chi.NewRouter()
	// creating file server for static files
	dir := http.Dir("./static")
	fs := http.FileServer(dir)

	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFound)
		components.Hello("Aman").Render(context.Background(), w)
	})
	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAlreadyReported)
		// FIX: Need to add html fileserver or serve files with templ (inefficient)
	})

	// Routes for User accounts
	r.Mount("/user/", UserRouter())

	// Routes for Manager accounts
	r.Mount("/manager", ManagerRouter())

	// Routes for admin accounts
	r.Mount("/admin", AdminRouter())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("route does not exist"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method is not valid"))
	})

	http.ListenAndServe(os.Getenv("PORT"), r)
}

// admin router configuration

// adminOnly midddleware that restricts access to just administrators
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("user.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// adminRouter is for routes for all admins
func AdminRouter() http.Handler {
	r := chi.NewRouter()
	// TODO: need to add authentication middleware
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello admin"))
	})
	return r
}

// ManagerOnly midddleware that restricts access to just administrators
func ManagerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isManager, ok := r.Context().Value("user.manager").(bool)
		if !ok || !isManager {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ManagerRouter is for routes for all managers
func ManagerRouter() http.Handler {
	r := chi.NewRouter()
	// TODO: need to add authentication middleware
	r.Use(ManagerOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Manager"))
	})
	return r
}

// UserOnly midddleware that restricts access to just administrators
func UserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isUser, ok := r.Context().Value("user.user").(bool)
		if !ok || !isUser {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UserRouter is for routes for all managers
func UserRouter() http.Handler {
	r := chi.NewRouter()
	// TODO: need to add authentication middleware
	r.Use(UserOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello User"))
	})
	return r
}
