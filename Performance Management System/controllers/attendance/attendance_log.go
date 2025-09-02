package attendance

import (
	"context"
	"net/http"
	"time"

	"performanceManagement/config"
	"performanceManagement/models/attendance"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AttendanceRequest struct {
	EmployeeID int        `json:"employee_id" binding:"required"`
	TargetDate *time.Time `json:"target_date"`
}

func GetAttendanceByEmpAndDate(c *gin.Context) {
	empIDStr := c.Param("employee_id")
	dateStr := c.Param("date")

	empID, err := strconv.Atoi(empIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing date"})
		return
	}

	rows, err := config.DB.Query(context.Background(), "SELECT * FROM fn_get_attendance_by_emp_and_date($1, $2)", empID, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var records []attendance.AttendanceRecordResponse
	for rows.Next() {
		var r attendance.AttendanceRecordResponse
		err := rows.Scan(&r.AttendanceID, &r.EmployeeID, &r.EmployeeName, &r.DepartmentName, &r.Date, &r.TimeIn, &r.TimeOut, &r.TotalHours, &r.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error"})
			return
		}
		records = append(records, r)
	}

	c.JSON(http.StatusOK, records)
}
