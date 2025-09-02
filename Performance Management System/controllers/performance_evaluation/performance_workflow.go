package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddPerformanceWorkflow(c *gin.Context) {
	var PW performance_evaluation.PerformanceWorkflow

	if err := c.ShouldBindJSON(&PW); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.InsertPeroformanceWorkflow(PW); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance Workflow Insertion failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Performance Workflow inserted successfully"})
}

func UPDATEPerformanceWorkflow(c *gin.Context) {
	idParam := c.Param("performance_workflow_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "invalid id"})
		return
	}

	var PW performance_evaluation.PerformanceWorkflow
	if err := c.ShouldBindJSON(&PW); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	PW.PerformanceWorkflowID = id
	if err := performance_evaluation.UpdatePerformanceWorkflow(PW); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance Workflow updation failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Performance Workflow updated successfully"})
}

func ListPerformanceWorkflow(c *gin.Context) {
	idParam := c.Param("performance_workflow_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	workflow, err := performance_evaluation.GetPerformanceWorkflow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to fetch the performance workflow": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, workflow)
}
