package controllers

import (
	"context"
	"fmt"
	"net/http"
	"performanceManagement/config"
	"performanceManagement/models"
	"performanceManagement/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	EmployeeID int    `json:"employee_id" binding:"required"`
	Role       string `json:"role"`
}

// RegisterHandler is not being used currently, we are using SignupHandler instead
// but it can be used for a different registration flow if needed.
func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		if req.Role == "" {
			req.Role = "employee"
		}
		if req.Role != "employee" && req.Role != "hr" && req.Role != "manager" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}

		var exists int
		err := config.DB.QueryRow(c, `SELECT 1 FROM users WHERE email = $1`, req.Email).Scan(&exists)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		_, err = config.DB.Exec(
			c.Request.Context(),
			`CALL sp_insert_user($1, $2, $3, $4)`,
			req.EmployeeID,
			req.Email,
			string(hashedPassword),
			req.Role,
		)

		if err != nil {
			fmt.Println("Error inserting the user: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var userID int
	var employeeID int
	var email string
	var hashedPassword string
	var role string

	query := `
		SELECT user_id, email, employee_id, password, role
		FROM users
		WHERE email = $1
	`

	err := config.DB.QueryRow(c, query, req.Email).Scan(&userID, &email, &employeeID, &hashedPassword, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(userID, employeeID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":       token,
		"role":        role,
		"user_id":     userID,
		"employee_id": employeeID,
	})
}

func SignupHandler(c *gin.Context) {
	var input models.SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var userID int
	err = config.DB.QueryRow(ctx, `
        INSERT INTO users (email, password, role)
        VALUES ($1, $2, $3)
        RETURNING user_id
    `, input.Email, string(hashedPassword), input.Role).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists or user insert failed"})
		return
	}

	// 2. Create blank employee and get ID
	var employeeID int
	err = config.DB.QueryRow(ctx, `
        INSERT INTO employees DEFAULT VALUES RETURNING employee_id
    `).Scan(&employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee record"})
		return
	}

	// 3. Update user to link employee_id
	_, err = config.DB.Exec(ctx, `
        UPDATE users SET employee_id = $1 WHERE user_id = $2
    `, employeeID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link employee_id to user"})
		return
	}

	// 4. Generate token
	token, err := utils.GenerateToken(userID, employeeID, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Signup successful",
		"token":       token,
		"role":        input.Role,
		"user_id":     userID,
		"employee_id": employeeID,
	})
}
