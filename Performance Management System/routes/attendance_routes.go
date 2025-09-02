package routes

import (
	"performanceManagement/controllers/attendance"
	"performanceManagement/middleware"

	"github.com/gin-gonic/gin"
)

func AttendanceRoutes(r *gin.Engine) {
	attendanceGroup := r.Group("/attendance")
	attendanceGroup.Use(middleware.JWTMiddleware())
	attendanceGroup.Use(middleware.PermissionMiddleware())
	{
		// IF hr directly want to post a leave application
		// This post api was used before the manager approval flow was implemented
		attendanceGroup.POST("/leave-applications", attendance.InsertLeaveApplication)
		// HR gets the list of leave applications
		attendanceGroup.GET("/leave-applications/:id", attendance.GetLeaveApplicationByID)

		//This endpoint is not for frontend. It is machine to machine endpoint because of use of webhook
		attendanceGroup.POST("/punch-logs", attendance.PunchHandler)

		// Employee apply for leave
		attendanceGroup.POST("/leave/apply", attendance.ApplyLeave)
		attendanceGroup.GET("/leave/manager/:manager_id", attendance.GetLeavesForManager)
		attendanceGroup.PUT("/leave/manager/approve/:leave_id", attendance.ManagerApproveLeave)
		attendanceGroup.GET("/leave/hr/pending", attendance.GetLeavesForHR)
		attendanceGroup.PUT("/leave/hr/approve/:leave_id", attendance.HRApproveLeave)
		// Attendance Records
		attendanceGroup.GET("/records/:employee_id/:date", attendance.GetAttendanceByEmpAndDate)
		// absence tracking
		attendanceGroup.GET("/absence-tracking/:date", attendance.GetAbsenceTracking)
		// Overtime
		attendanceGroup.POST("/overtime/calculate/:employee_id/:date", attendance.CalculateOvertime)
		attendanceGroup.GET("/overtime/all", attendance.GetAllCalculatedOvertime)

		// Monthly summmary
		attendanceGroup.GET("/:id/monthly-summary/:month/:year", attendance.GetMonthlySummary)
		attendanceGroup.POST("/monthly-summary", attendance.GenerateMonthlySummary)
		// Attedance Alerts
		attendanceGroup.GET("/alerts", attendance.GetAttendanceAlerts)
		attendanceGroup.POST("/alerts", attendance.CreateAttendanceAlert)
		attendanceGroup.PUT("/alerts/:alert_id", attendance.UpdateAttendanceAlert)

	}
}
