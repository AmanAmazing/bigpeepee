package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"purchaseOrderSystem/auth"
	"purchaseOrderSystem/models"
	"purchaseOrderSystem/utils"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{
		db: db}
}

func (s *UserService) Signup(email, username, password string) error {

	hashedpassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(context.Background(), "insert into users (email,username,password) values ($1,$2,$3)", email, username, hashedpassword)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			log.Printf("username %s has already been registered", username)
		}
		return err
	}
	return nil
}

func (s *UserService) Login(username, password string) (string, string, error) {
	// finding the user
	var user models.User
	err := s.db.QueryRow(context.Background(), "SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("invalid username or password")
		}
		return "", "", fmt.Errorf("failed to retrieve user: %v", err)
	}
	//
	match := utils.CheckPasswordMatch(user.Password, password)
	if !match {
		return "", "", fmt.Errorf("invalid password")
	}
	// Query the user_roles table to get the user's role
	var roleID int
	err = s.db.QueryRow(context.Background(), "SELECT role_id FROM user_roles WHERE user_id = $1", user.ID).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("user role not found")
		}
		return "", "", fmt.Errorf("failed to retrieve user role: %v", err)
	}

	// Query the roles table to get the role name
	err = s.db.QueryRow(context.Background(), "SELECT name FROM roles WHERE id = $1", roleID).Scan(&user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("role not found")
		}
		return "", "", fmt.Errorf("failed to retrieve role name: %v", err)
	}

	// Getting the user's department
	var departmentID int
	err = s.db.QueryRow(context.Background(), "SELECT department_id FROM user_departments WHERE user_id = $1", user.ID).Scan(&departmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("user department not found")
		}
		return "", "", fmt.Errorf("failed to retrieve user department: %v", err)
	}
	// Query the departments table to get the department name
	err = s.db.QueryRow(context.Background(), "SELECT name FROM departments WHERE id = $1", departmentID).Scan(&user.Department)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("department not found")
		}
		return "", "", fmt.Errorf("failed to retrieve department name: %v", err)
	}

	// Generate a JWT token
	claims := map[string]interface{}{
		"userId":     user.ID,
		"role":       user.Role,
		"department": user.Department,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Sign the token with a secret key
	tokenString, err := auth.GenerateJWT(claims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %v", err)
	}
	return tokenString, user.Role, nil
}

type Supplier struct {
	ID   int
	Name string
}
type Nominal struct {
	ID   int
	Name string
}

type Product struct {
	ID   int
	Name string
}

func (s *UserService) GetSuppliers() ([]Supplier, error) {
	rows, err := s.db.Query(context.Background(), "SELECT id, name FROM suppliers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []Supplier
	for rows.Next() {
		var supplier Supplier
		err := rows.Scan(&supplier.ID, &supplier.Name)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil

}

func (s *UserService) GetNominals() ([]Nominal, error) {
	rows, err := s.db.Query(context.Background(), "SELECT id, name FROM nominals")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nominals []Nominal
	for rows.Next() {
		var nominal Nominal
		err := rows.Scan(&nominal.ID, &nominal.Name)
		if err != nil {
			return nil, err
		}
		nominals = append(nominals, nominal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nominals, nil
}

func (s *UserService) GetProducts() ([]Product, error) {
	rows, err := s.db.Query(context.Background(), "SELECT id, name FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
