package controllers

import (
	"context"
	"net/http"
	"performanceManagement/config"

	"github.com/gin-gonic/gin"
)

type RolePermissionRequest struct {
	RoleID   int    `json:"role_id"`
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

type UserPermissionRequest struct {
	UserID         int    `json:"user_id"`
	PermissionName string `json:"permission_name"`
}

// Assign permission to a ROLE (affects all users of that role)
func AssignRolePermission(c *gin.Context) {
	var req RolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(
		context.Background(),
		`
	INSERT INTO api_permissions (role_id, method, endpoint)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	`,
		req.RoleID, req.Method, req.Endpoint,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to role"})
}

// Assign permission to a specific USER
func AssignUserPermission(c *gin.Context) {
	var req UserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find permission_id
	var permID int
	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id FROM permissions WHERE name = $1`,
		req.PermissionName,
	).Scan(&permID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission name"})
		return
	}

	_, err = config.DB.Exec(
		context.Background(),
		`INSERT INTO user_permissions (user_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`,
		req.UserID, permID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to user"})
}

type CreatePermissionRequest struct {
	Name string `json:"name"`
}

func CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission name is required"})
		return
	}

	_, err := config.DB.Exec(
		context.Background(),
		`INSERT INTO permissions (name) VALUES ($1) ON CONFLICT DO NOTHING`,
		req.Name,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission created successfully"})
}
