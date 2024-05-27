package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"purchaseOrderSystem/auth"
	"purchaseOrderSystem/models"
	"purchaseOrderSystem/utils"
	"strconv"
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

func (s *UserService) SubmitPurchaseOrder(userID int, department, priority string, item_count int, formData url.Values) error {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to start a transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	var departmentID int
	err = tx.QueryRow(context.Background(), `
		SELECT id FROM departments WHERE name = $1
	`, department).Scan(&departmentID)
	if err != nil {
		return fmt.Errorf("failed to get department ID: %v", err)
	}

	title := formData.Get("title")
	description := formData.Get("description")

	var purchaseOrderID int
	err = tx.QueryRow(context.Background(), `
		INSERT INTO purchase_orders (user_id, department_id,title,description,priority) 
		VALUES ($1,$2,$3,$4,$5) 
		RETURNING id`, userID, departmentID, title, description, priority).Scan(&purchaseOrderID)
	if err != nil {
		return fmt.Errorf("failed to insert purchase order 1: %v", err)
	}

	for i := 1; i <= item_count; i++ {
		itemName := formData.Get(fmt.Sprintf("name%d", i))
		supplier := formData.Get(fmt.Sprintf("supplier%d", i))
		nominal := formData.Get(fmt.Sprintf("nominal%d", i))
		product := formData.Get(fmt.Sprintf("product%d", i))
		unitPriceStr := formData.Get(fmt.Sprintf("unit_price%d", i))
		quantityStr := formData.Get(fmt.Sprintf("quantity%d", i))
		link := formData.Get(fmt.Sprintf("link%d", i))

		unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
		if err != nil {
			return fmt.Errorf("Invalid unit price for item %d: %v", i, err)
		}
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return fmt.Errorf("Invalid quantity for item %d: %v", i, err)
		}
		totalPrice := unitPrice * float64(quantity)

		_, err = tx.Exec(context.Background(), `
			INSERT INTO purchase_order_items (purchase_order_id,item_name,supplier,nominal,product,unit_price,quantity,total_price,link)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			`, purchaseOrderID, itemName, supplier, nominal, product, unitPrice, quantity, totalPrice, link)
		if err != nil {
			return fmt.Errorf("failed to insert purchase order 2 item: %v", err)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// func (s *UserService) ProcessLoginForm(username, password string) (string, string, error) {
//
// }
