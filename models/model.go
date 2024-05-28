package models

import "time"

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Department string `json:"department"`
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

type PurchaseOrder struct {
	ID           int
	UserID       int
	DepartmentID int
	Title        string
	Description  string
	Status       string
	Priority     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type PurchaseOrderItem struct {
	ID              int
	PurchaseOrderID int
	ItemName        string
	Supplier        string
	Nominal         string
	Product         string
	Quantity        int
	UnitPrice       float64
	TotalPrice      float64
	Link            string
	Status          string
	ApproverID      *int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

// type FormItem struct {
// 	Name      string  `form:"name"`
// 	Supplier  string  `form:"supplier"`
// 	Nominal   string  `form:"nominal"`
// 	Product   string  `form:"product"`
// 	UnitPrice float64 `form:"unit_price"`
// 	Quantity  int     `form:"quantity"`
// 	Link      string  `form:"link"`
// }
