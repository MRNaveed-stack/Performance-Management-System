package routes

import (
	"performanceManagement/controllers"
	"performanceManagement/middleware"

	"github.com/gin-gonic/gin"
)

func EmployeeRoutes(r *gin.Engine) {
	auth := r.Group("")
	auth.Use(middleware.JWTMiddleware())
	auth.Use(middleware.PermissionMiddleware())
	{
		emp := auth.Group("/employees")
		{
			emp.GET("/basic", controllers.GetAllEmployeeBasicInfo)
			emp.GET("/basic/:id", controllers.GetEmployeeBasicInfoByIDHandler)
			emp.POST("/register", controllers.CreateEmployeeBasic)
			emp.PUT("/basic/:id", controllers.UpdateEmployeeBasicInfo)

			// Employment Details
			emp.GET("/:id/employment-details", controllers.GetEmploymentDetails)
			emp.POST("/:id/employment-details", controllers.SaveEmploymentDetailsHandler)
			emp.PUT("/:id/employment-details", controllers.UpdateEmploymentDetails)

			// Salary
			emp.GET("/:id/salary", controllers.GetSalaryInfoByEmployeeID)
			emp.POST("/:id/salary", controllers.CreateSalaryInfo)
			emp.PUT("/salary/:salary_id", controllers.UpdateSalaryInfo)
			emp.DELETE("/salary/:salary_id", controllers.DeleteSalaryInfo)

			// Education
			emp.GET("/:id/education", controllers.GetEducation)
			emp.POST("/:id/education", controllers.AddEducation)
			emp.DELETE("/education/:education_id", controllers.DeleteEducation)

			// Certifications
			emp.GET("/:id/certifications", controllers.GetCertifications)
			emp.POST("/:id/certifications", controllers.AddCertification)
			emp.DELETE("/certifications/:cert_id", controllers.DeleteCertification)

			// Documents
			emp.GET("/employee/documents/:doc_id", controllers.GetSingleEmployeeDocument)
			emp.POST("/employee/:id/documents", controllers.UploadEmployeeDocument)
			emp.DELETE("/employee/documents/:doc_id", controllers.DeleteEmployeeDocument)

			// Skills
			emp.GET("/:id/skills", controllers.GetEmployeeSkillsHandler)
			emp.POST("/:id/skills", controllers.AddEmployeeSkillHandler)
			emp.POST("/skills", controllers.CreateSkillHandler)
			emp.DELETE("/:id/skills/:skill_id", controllers.DeleteEmployeeSkillHandler)

			// Reporting Structure
			emp.GET("/:id/reporting-structure", controllers.GetFullReportingStructureByID)
			emp.PUT("/:id/reporting-structure", controllers.UpdateReportingStructure)

			// Employment History
			emp.GET("/employee/history", controllers.GetAllEmploymentHistoryHandler)
			emp.GET("/employee/:id/history", controllers.GetEmployeeHistoryByIDHandler)
			emp.POST("/employee/:id/history", controllers.InsertEmploymentHistoryHandler)

		}
	}
}

func LookupRoutes(r *gin.Engine) {
	lookup := r.Group("/lookups")

	// Public endpoints
	lookup.GET("/genders", controllers.GetAllGenders)
	lookup.GET("/marital-status", controllers.GetAllMaritalStatuses)

	// Department info requires authentication
	authLookup := lookup.Group("")
	authLookup.Use(middleware.JWTMiddleware())
	authLookup.Use(middleware.PermissionMiddleware())
	{
		authLookup.GET("/departments", controllers.GetAllDepartments)
	}
}
