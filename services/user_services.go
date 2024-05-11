package services

import (
	"context"
	"log"
	"purchaseOrderSystem/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Signup(email, username, password string) error {

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(context.Background(), "INSERT INTO users (email,username,password) VALUES ($1,$2,$3)", email, username, hashedPassword)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			log.Printf("username %s has already been registered", username)
		}
		return err
	}
	return nil
}
