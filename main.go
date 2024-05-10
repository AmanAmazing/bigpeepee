package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"purchaseOrderSystem/components"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// FIX: Need to fix templ support for neovim
var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
	// initiating database connection
	// TODO: Need to separate this out into another module
	dsn := fmt.Sprintf("%s://%s:%s@localhost:%s/%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db_pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db_pool.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(jwtauth.Verifier(tokenAuth))

	// creating file server for static files
	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public unathenticated routes
	r.Mount("/", PublicRouter())

	// Routes for User accounts
	r.Mount("/user", UserRouter())

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

// This router is for all urls that are publically accessible
func PublicRouter() http.Handler {
	r := chi.NewRouter()
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
		// TODO: Link with signup services
	})
	return r
}

// admin router configuration

// adminOnly midddleware that restricts access to just administrators
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		if claims["role"] != "admin" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// adminRouter is for routes for all admins
func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello admin"))
	})
	return r
}

// ManagerOnly midddleware that restricts access to just administrators
func ManagerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		if claims["role"] != "manager" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ManagerRouter is for routes for all managers
func ManagerRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(ManagerOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Manager"))
	})
	return r
}

// UserOnly midddleware that restricts access to just administrators
func UserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		if claims["role"] != "user" {
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

func GenerateJWT(claims map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
