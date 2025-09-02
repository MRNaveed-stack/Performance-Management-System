package attendance

import (
	"log"
	"net/http"
	"performanceManagement/config"
	"performanceManagement/models/attendance"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func CalculateOvertime(c *gin.Context) {
	employeeID := c.Param("employee_id")
	overtimeDate := c.Param("date")

	// Request body from HR
	var input struct {
		WorkingHours float64 `json:"working_hours" binding:"required"`
		HourlyRate   float64 `json:"hourly_rate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", overtimeDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// 1. Get department_id from employees table
	var departmentID int
	err = config.DB.QueryRow(
		c,
		`SELECT employee_department_id FROM employees WHERE employee_id = $1`,
		employeeID,
	).Scan(&departmentID)
	if err != nil {
		log.Printf("Error fetching department_id for employee_id=%s: %v", employeeID, err)
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department_id"})
		return
	}

	// 2. Get earliest punch_in and latest punch_out for the day
	var timeIn, timeOut time.Time
	err = config.DB.QueryRow(
		c,
		`SELECT 
			MIN(pl.punch_date_time) FILTER (WHERE pt.punch_type_name = 'punch_in') AS time_in,
			MAX(pl.punch_date_time) FILTER (WHERE pt.punch_type_name = 'punch_out') AS time_out
		FROM punch_logs pl
		JOIN punch_type pt ON pl.punch_type_id = pt.punch_type_id
		WHERE pl.employee_id = $1 AND DATE(pl.punch_date_time) = $2`,
		employeeID, date,
	).Scan(&timeIn, &timeOut)
	if err != nil || timeIn.IsZero() || timeOut.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing punch_in or punch_out for this date"})
		return
	}

	// 3. Calculate total and overtime hours
	totalHours := timeOut.Sub(timeIn).Hours()
	overtimeHours := totalHours - input.WorkingHours
	if overtimeHours < 0 {
		overtimeHours = 0
	}

	// 4. Calculate overtime payment
	overtimePayment := overtimeHours * input.HourlyRate

	// 5. Save to DB
	err = attendance.SaveOvertimeRecord(
		mustAtoi(employeeID),
		departmentID,
		date,
		timeIn,
		timeOut,
		totalHours,
		overtimeHours,
		input.HourlyRate,
		overtimePayment,
	)
	if err != nil {
		log.Printf("Error saving overtime record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save overtime record"})
		return
	}

	// 6. Respond with calculated details
	c.JSON(http.StatusOK, gin.H{
		"employee_id":            employeeID,
		"employee_department_id": departmentID,
		"overtime_date":          overtimeDate,
		"time_in":                timeIn.Format("15:04"),
		"time_out":               timeOut.Format("15:04"),
		"total_hours":            totalHours,
		"overtime_hours":         overtimeHours,
		"hourly_rate":            input.HourlyRate,
		"overtime_payment":       overtimePayment,
	})
}

// helper function to convert string to int
func mustAtoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func GetAllCalculatedOvertime(c *gin.Context) {
	rows, err := config.DB.Query(c, `
		SELECT 
			employee_id,
			department_id,
			overtime_date::text,
			to_char(time_in, 'HH24:MI') AS time_in,
			to_char(time_out, 'HH24:MI') AS time_out,
			total_hours,
			overtime_hours,
			hourly_rate,
			overtime_payment
		FROM overtime_records
		ORDER BY overtime_date DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch overtime records"})
		return
	}
	defer rows.Close()

	var records []map[string]interface{}

	for rows.Next() {
		var (
			employeeID      int
			departmentID    int
			overtimeDate    string
			timeIn          string
			timeOut         string
			totalHours      float64
			overtimeHours   float64
			hourlyRate      float64
			overtimePayment float64
		)

		if err := rows.Scan(
			&employeeID,
			&departmentID,
			&overtimeDate,
			&timeIn,
			&timeOut,
			&totalHours,
			&overtimeHours,
			&hourlyRate,
			&overtimePayment,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading overtime record"})
			return
		}

		records = append(records, map[string]interface{}{
			"employee_id":      employeeID,
			"department_id":    departmentID,
			"overtime_date":    overtimeDate,
			"time_in":          timeIn,
			"time_out":         timeOut,
			"total_hours":      totalHours,
			"overtime_hours":   overtimeHours,
			"hourly_rate":      hourlyRate,
			"overtime_payment": overtimePayment,
		})
	}

	c.JSON(http.StatusOK, records)
}
