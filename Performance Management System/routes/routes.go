package routes

import (
	"performanceManagement/controllers"
	"performanceManagement/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/login", controllers.LoginHandler)
	// r.POST("/register", controllers.RegisterHandler())
	r.POST("/auth/signup", controllers.SignupHandler)
	auth := r.Group("/employees")
	auth.Use(middleware.AuthMiddleware())

}
