package main

import (
	"log"
	"net/http"
	"os"
	"purchaseOrderSystem/database"
	"purchaseOrderSystem/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
}

func main() {
	// initiating database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// insert test data into the database
	// database.TestDB() // FIX: not best practice. Need to return err from this

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// creating file server for static files
	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public unathenticated routes
	r.Mount("/", routes.PublicRouter(db))
	// Routes for User accounts
	r.Mount("/user", routes.UserRouter(db))

	// Routes for admin accounts
	r.Mount("/admin", routes.AdminRouter())

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
