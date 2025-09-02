package attendance

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"performanceManagement/config"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAttendanceAlerts(c *gin.Context) {
	deptIDStr := c.Query("department_id")
	dateStr := c.Query("date")
	var deptID interface{}
	var alertDate interface{}

	if deptIDStr != "" {
		if id, err := strconv.Atoi(deptIDStr); err == nil {
			deptID = id
		}
	} else {
		deptID = nil
	}

	if dateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			alertDate = parsedDate
		}
	} else {
		alertDate = nil
	}

	rows, err := config.DB.Query(
		context.Background(),
		`SELECT alert_id, employee_id, employee_name, department_name, date, remarks, category
     FROM fn_list_attendance_alerts($1, $2)`,
		deptID, alertDate,
	)

	if err != nil {
		log.Printf("Error fetching alerts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var alerts []map[string]interface{}
	for rows.Next() {
		var alertID, empID int
		var name, deptName, remark, category sql.NullString
		var date time.Time

		err := rows.Scan(&alertID, &empID, &name, &deptName, &date, &remark, &category)
		if err != nil {
			log.Printf("scan error: %v", err)
			continue
		}

		alert := gin.H{
			"alert_id":        alertID,
			"employee_id":     empID,
			"employee_name":   name.String,
			"department_name": deptName.String,
			"date":            date.Format("2006-01-02"),
			"remarks":         remark.String,
			"category":        category.String,
		}
		alerts = append(alerts, alert)
	}

	c.JSON(http.StatusOK, alerts)
}

func CreateAttendanceAlert(c *gin.Context) {
	var req struct {
		EmployeeID   int    `json:"employee_id"`
		DepartmentID int    `json:"department_id"`
		Date         string `json:"date"`
		Remark       string `json:"remarks"`
		Category     string `json:"category"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dateParsed, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	_, err = config.DB.Exec(context.Background(), `
		CALL sp_insert_attendance_alert($1, $2, $3, $4, $5)
	`, req.EmployeeID, req.DepartmentID, dateParsed, req.Remark, req.Category)

	if err != nil {
		fmt.Println("DB insert error:", err) // log actual cause
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance alert created"})
}

func UpdateAttendanceAlert(c *gin.Context) {
	alertIDStr := c.Param("alert_id")
	alertID, err := strconv.Atoi(alertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	var req struct {
		Remark       string `json:"remarks"`
		Category     string `json:"category"`
		DepartmentID int    `json:"department_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = config.DB.Exec(context.Background(), `
		CALL sp_update_attendance_alert($1, $2, $3, $4)
	`, alertID, req.Remark, req.Category, req.DepartmentID)

	if err != nil {
		log.Printf("Failed to update alert: alert_id=%d, remarks=%q, category=%q, department_id=%d, error=%v",
			alertID, req.Remark, req.Category, req.DepartmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance alert updated"})
}
