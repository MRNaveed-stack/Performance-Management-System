package routes

import (
	"performanceManagement/controllers/performance_evaluation"
	"performanceManagement/middleware"

	"github.com/gin-gonic/gin"
)

func PerformanceEvaluationRoutes(r *gin.Engine) {
	attendanceGroup := r.Group("/evaluation")
	attendanceGroup.Use(middleware.JWTMiddleware())
	attendanceGroup.Use(middleware.PermissionMiddleware())
	{

		// SETUP PHASE
		r.POST("/category-create", performance_evaluation.AddKpiCategory)
		r.PUT("/category-update/:kpi_category_id", performance_evaluation.UpdateKpiCategory)
		r.GET("/category-list", performance_evaluation.GetKpiCategory)

		r.POST("/create-kpi", performance_evaluation.AddKpi)
		r.PUT("/update-kpi/:kpi_id", performance_evaluation.UpdateKpi)
		r.GET("/get-kpi", performance_evaluation.UpdateKpiScoringCriteria)

		r.POST("/create-kpiscroingcriteria", performance_evaluation.AddKpiScoringCriteria)
		r.PUT("/update-kpiscoringcriteria/:kpi_scoring_criteria_id", performance_evaluation.UpdateKpiScoringCriteria)
		r.GET("/:id/get-kpiscoringcriteria", performance_evaluation.ListKpiScoringCriteria)

		r.POST("/Add-kpi-target", performance_evaluation.AddKpiTarget)
		r.PUT("/Update-kpi-target/:kpi_target_id", performance_evaluation.UPDATEKpiTarget)
		r.GET("/Get-kpi-target", performance_evaluation.ListKpiTarget)

		r.POST("/Assign-kpi", performance_evaluation.AddEmployeeKpiAssignment)
		r.PUT("/update-assigned-kpi/:employee_kpi_assignment_id", performance_evaluation.UPDATEEmployeeKpiAssignment)
		r.GET("/Get-assigned-kpi/:employee_id", performance_evaluation.ListEmployeeKpiAssigment)

		r.POST("/Add-performance-cycle", performance_evaluation.AddPerformanceCycle)
		r.PUT("/Update-performance-cycle/:performance_cycle_id", performance_evaluation.UPDATEPerformanceCycle)
		r.GET("/Get-performance-cycle", performance_evaluation.ListPerformanceCycle)

		r.POST("/Add-self-appraisal", performance_evaluation.AddPerformanceCycle)
		r.PUT("/Update-self-appraisal", performance_evaluation.UPDATEPerformanceCycle)
		r.GET("/Get-self-appraisal", performance_evaluation.ListPerformanceCycle)
		r.GET("/Get-self-appraisal-id/:self_appraisal_id", performance_evaluation.ListSelfAppraisalByID)

		// REVIEW PHASE
		r.POST("/Add-performance-review", performance_evaluation.AddPerformanceReview)
		r.PUT("/Update-performance-review/:performance_review_id", performance_evaluation.UPDATEPerformanceReview)
		r.GET("/Get-performance-review", performance_evaluation.ListPerformanceReview)
		r.GET("/Get-performance-review/:performance_review_id", performance_evaluation.ListPerformanceReviewByID)

		r.POST("/Add-performance-workflow", performance_evaluation.AddPerformanceWorkflow)
		r.PUT("/Update-performance-workflow/:performance_workflow_id", performance_evaluation.UPDATEPerformanceWorkflow)
		r.GET("/Get-performance-workflow/:performance_workflow_id", performance_evaluation.ListPerformanceWorkflow)

		r.POST("/Add-performance-approval", performance_evaluation.AddPerformanceApproval)
		r.PUT("/Update-performance-approval/:approval_id", performance_evaluation.UPDATEPerformanceApproval)
		r.GET("/Get-performance-approval/:approval_id", performance_evaluation.ListPerformanceApproval)

		r.POST("/Performance-document", performance_evaluation.AddPerformanceDocument)
		r.GET("/Performance-document", performance_evaluation.ListPerformanceDocuments)
		r.GET("/Performance-document/download/:document_id", performance_evaluation.DownloadPerformanceDocument)
		r.DELETE("/Delete-document/:document_id", performance_evaluation.DeletePerformanceDocument)
		r.Static("/uploads", "./uploads/")

		r.POST("/Add-performance-attendance", performance_evaluation.AddPerformanceAttendance)
		r.PUT("/Update-performance-attendance/:attendance_id", performance_evaluation.UpdatePerformanceAttendance)
		r.GET("/Get-performance-attendance/:employee_id", performance_evaluation.ListPerformanceAttendance)

		// Finalization phase
	}

}
