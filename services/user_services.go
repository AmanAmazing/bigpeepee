package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"purchaseOrderSystem/auth"
	"purchaseOrderSystem/models"
	"purchaseOrderSystem/utils"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
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

func (s *UserService) GetSuppliers() ([]models.Supplier, error) {
	rows, err := s.db.Query(context.Background(), "SELECT id, name FROM suppliers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var supplier models.Supplier
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

func (s *UserService) GetNominals() ([]models.Nominal, error) {
	rows, err := s.db.Query(context.Background(), "SELECT id, name FROM nominals")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nominals []models.Nominal
	for rows.Next() {
		var nominal models.Nominal
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

func (s *UserService) GetProducts() ([]models.Product, error) {
	rows, err := s.db.Query(context.Background(), "SELECT id, name FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
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

func (s *UserService) SubmitPurchaseOrder(userID int, department, priority, role string, item_count int, formData url.Values) error {
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

		var status string

		switch role {
		case "user":
			status = "pending"
		case "manager":
			status = "manager_approved"
		case "admin": // TODO: Not really sure I should let admin users submit POs might cancell this later on
			status = "pending"
		}

		_, err = tx.Exec(context.Background(), `
			INSERT INTO purchase_order_items (purchase_order_id,item_name,supplier,nominal,product,unit_price,quantity,total_price,link, status)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			`, purchaseOrderID, itemName, supplier, nominal, product, unitPrice, quantity, totalPrice, link, status)
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

// type FormItem struct {
// Queries the database for all purchase orders created by the supplied userId
func (s *UserService) GetPurchaseOrdersByUserID(userId int) ([]models.PurchaseOrder, error) {
	query := fmt.Sprintf("SELECT * FROM purchase_orders WHERE user_id=%v LIMIT 10", userId)
	fmt.Println("Query is: ", query)

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Unable to query for purchase orders: %w", err)
	}
	defer rows.Close()

	var purchaseOrders []models.PurchaseOrder
	for rows.Next() {
		po := models.PurchaseOrder{}
		err := rows.Scan(&po.ID, &po.UserID, &po.DepartmentID, &po.Title, &po.Description, &po.Status, &po.Priority, &po.CreatedAt, &po.UpdatedAt, &po.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		purchaseOrders = append(purchaseOrders, po)
	}
	fmt.Println(len(purchaseOrders))
	return purchaseOrders, nil
}
func (s *UserService) GetPurchaseOrderById(userID int, poID string) (map[models.PurchaseOrder][]models.PurchaseOrderItem, error) {
	purchaseOrder := make(map[models.PurchaseOrder][]models.PurchaseOrderItem)

	// Start a transaction
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() { // Rollback if anything goes wrong
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// 1. Retrieve the purchase order
	poQuery := `
        SELECT id, user_id, department_id, title, description, status, priority, created_at, updated_at
        FROM purchase_orders
        WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
    `
	var po models.PurchaseOrder
	err = tx.QueryRow(context.Background(), poQuery, poID, userID).Scan(
		&po.ID, &po.UserID, &po.DepartmentID, &po.Title, &po.Description, &po.Status, &po.Priority, &po.CreatedAt, &po.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("purchase order not found for this user")
		}
		return nil, fmt.Errorf("error fetching purchase order: %v", err)
	}

	// 2. Retrieve the associated purchase order items
	itemsQuery := `
        SELECT id, purchase_order_id, item_name, supplier, nominal, product, quantity, unit_price, total_price, link, status, approver_id, created_at, updated_at
        FROM purchase_order_items
        WHERE purchase_order_id = $1 AND deleted_at IS NULL
    `
	rows, err := tx.Query(context.Background(), itemsQuery, poID)
	if err != nil {
		return nil, fmt.Errorf("error fetching purchase order items: %v", err)
	}
	defer rows.Close()

	var items []models.PurchaseOrderItem
	for rows.Next() {
		var item models.PurchaseOrderItem
		err := rows.Scan(
			&item.ID, &item.PurchaseOrderID, &item.ItemName, &item.Supplier, &item.Nominal, &item.Product, &item.Quantity, &item.UnitPrice, &item.TotalPrice, &item.Link, &item.Status, &item.ApproverID, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning purchase order item: %v", err)
		}
		items = append(items, item)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating purchase order items: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	// Add the purchase order and its items to the result map
	purchaseOrder[po] = items
	return purchaseOrder, nil
}

func (s *UserService) GetPurchaseOrderByIdWithoutItems(poID string) (models.PurchaseOrder, error) {
	po := models.PurchaseOrder{}

	// 1. Retrieve the purchase order
	poQuery := `
        SELECT id, user_id, department_id, title, description, status, priority, created_at, updated_at
        FROM purchase_orders
        WHERE id = $1
    `
	err := s.db.QueryRow(context.Background(), poQuery, poID).Scan(
		&po.ID, &po.UserID, &po.DepartmentID, &po.Title, &po.Description, &po.Status, &po.Priority, &po.CreatedAt, &po.UpdatedAt,
	)
	if err != nil {
		return po, fmt.Errorf("error fetching purchase order: %v", err)
	}

	// Add the purchase order and its items to the result map
	return po, nil
}

// returns empty struct with error if it occurs
func (s *UserService) PutPurchaseOrder(formData models.PurchaseOrder) (models.PurchaseOrder, error) {

	statement := `UPDATE purchase_orders
				  SET title = $1, description = $2, priority = $3, 
				  updated_at = NOW()
				  WHERE id = $4
		`
	_, err := s.db.Exec(context.Background(), statement, formData.Title, formData.Description, formData.Priority, formData.ID)
	if err != nil {
		return models.PurchaseOrder{}, err
	}
	// check if the form data is the same
	return formData, err
}

// func (s *UserService) ProcessLoginForm(username, password string) (string, string, error) {
//
// }
