package attendance

import (
	"context"
	"fmt"
	"net/http"
	"performanceManagement/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateMonthlySummary(c *gin.Context) {
	var req struct {
		EmployeeID int `json:"employee_id"`
		Month      int `json:"month"`
		Year       int `json:"year"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := config.DB.Exec(context.Background(),
		"SELECT fn_generate_monthly_summary($1, $2, $3)",
		req.EmployeeID, req.Month, req.Year)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monthly summary generated"})
}

func GetMonthlySummary(c *gin.Context) {
	employeeIDStr := c.Param("id")
	monthStr := c.Param("month")
	yearStr := c.Param("year")

	employeeID, err := strconv.Atoi(employeeIDStr)
	month, err2 := strconv.Atoi(monthStr)
	year, err3 := strconv.Atoi(yearStr)

	if err != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID, month, or year"})
		return
	}

	var summary struct {
		TotalWorkingDays int `json:"total_working_days"`
		DaysPresent      int `json:"days_present"`
		DaysAbsent       int `json:"days_absent"`
		LateComings      int `json:"late_comings"`
	}

	err = config.DB.QueryRow(
		context.Background(),
		"SELECT total_working_days, days_present, days_absent, late_comings FROM  fn_get_monthly_attendance_summary($1, $2, $3)",
		employeeID, month, year,
	).Scan(&summary.TotalWorkingDays, &summary.DaysPresent, &summary.DaysAbsent, &summary.LateComings)

	if err != nil {
		fmt.Println("DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
