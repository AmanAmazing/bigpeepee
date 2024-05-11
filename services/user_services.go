package services

import (
	"context"
	"errors"
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
			return errors.New("username has already been registered")
		}
		return err
	}
	return nil
}
