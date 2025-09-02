package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddEmployeeKpiAssignment(c *gin.Context) {
	var EKA performance_evaluation.EmployeeKpiAssignment

	if err := c.ShouldBindJSON(&EKA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.InsertEmployeeKpiAssignment(EKA); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Kpi assignment to employee failed": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "Kpi is assigned to employee"})
}

func UPDATEEmployeeKpiAssignment(c *gin.Context) {
	idParam := c.Param("employee_kpi_assignment_id")
	var EKA performance_evaluation.EmployeeKpiAssignment
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := c.ShouldBindJSON(&EKA); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	EKA.EmployeeKpiAssignmentID = id
	if err := performance_evaluation.UpdateEmployeeKpiAssignment(EKA); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Updation of kpi assignment failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Updation of kpi assignment to employee is done successfully"})

}

func ListEmployeeKpiAssigment(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpis, err := performance_evaluation.GetEmployeeKpiAssignment(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to fetch kpi assignment": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, kpis)
}
