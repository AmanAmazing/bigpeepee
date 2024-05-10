package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type userService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) Signup(email, username, password string) error {
	return nil
}
