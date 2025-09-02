package models

type SignupInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type User struct {
	UserID     int    `json:"user_id"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Role       string `json:"role"`
	EmployeeID *int   `json:"employee_id,omitempty"`
}
