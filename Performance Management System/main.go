package main

import (
	"log"
	"performanceManagement/config"

	"performanceManagement/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(" Error loading .env file")
	}

	config.ConnectDB()
	defer config.CloseDB()

	router := gin.Default()

	routes.AuthRoutes(router)
	for _, ri := range router.Routes() {
		log.Printf("%s %s\n", ri.Method, ri.Path)
	}

	routes.EmployeeRoutes(router)
	routes.LookupRoutes(router)
	routes.PermissionRoutes(router)
	routes.AttendanceRoutes(router)

	router.Run(":9090")
}
