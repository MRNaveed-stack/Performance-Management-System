package routes

import (
	"performanceManagement/controllers"
	"performanceManagement/middleware"

	"github.com/gin-gonic/gin"
)

func PermissionRoutes(r *gin.Engine) {
	auth := r.Group("/permissions")
	auth.Use(middleware.JWTMiddleware())
	auth.Use(middleware.PermissionMiddleware()) 
	{
		auth.POST("/assign-role", controllers.AssignRolePermission)
		auth.POST("/create", controllers.CreatePermission)

		auth.POST("/assign-user", controllers.AssignUserPermission)
	}
}
