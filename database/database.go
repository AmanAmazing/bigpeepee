package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("%s://%s:%s@localhost:%s/%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	return pgxpool.New(context.Background(), dsn)
}

func TestDB() {
	test_data := "test_data.sql"
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("%s://%s:%s@localhost:%s/%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	script, err := os.ReadFile(test_data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read SQL script: %v\n", err)
		os.Exit(1)
	}
	_, err = conn.Exec(context.Background(), string(script))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while executing SQL script: %v\n", err)
		os.Exit(1)
	}

	log.Println("test data inserted successfully")
}
