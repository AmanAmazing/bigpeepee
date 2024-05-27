package models

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Department string `json:"department"`
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
