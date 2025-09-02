package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Insert
func AddPerformanceAttendance(c *gin.Context) {
	var pa performance_evaluation.PerformanceAttendance
	if err := c.ShouldBindJSON(&pa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := performance_evaluation.InsertPerformanceAttendance(pa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendance inserted successfully", "attendance_id": id})
}

// Update
func UpdatePerformanceAttendance(c *gin.Context) {
	idParam := c.Param("attendance_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance_id"})
		return
	}

	var pa performance_evaluation.PerformanceAttendance
	if err := c.ShouldBindJSON(&pa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pa.AttendanceID = id

	if err := performance_evaluation.UpdatePerformanceAttendance(pa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully"})
}

// List by employee
func ListPerformanceAttendance(c *gin.Context) {
	employeeIDParam := c.Param("employee_id")
	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee_id"})
		return
	}

	records, err := performance_evaluation.ListPerformanceAttendance(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}
