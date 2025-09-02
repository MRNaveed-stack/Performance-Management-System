package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddPerformanceCycle(c *gin.Context) {
	var PC performance_evaluation.PerformanceCycle

	if err := c.ShouldBindJSON(&PC); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.InsertPerformanceCycle(PC); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance Cycle Insertion Failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Performance Cycle inserted successfully"})
}

func UPDATEPerformanceCycle(c *gin.Context) {
	idParam := c.Param("performance_cycle_id")
	var PC performance_evaluation.PerformanceCycle
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := c.ShouldBindJSON(&PC); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	PC.PerformanceCycleID = id
	if err := performance_evaluation.UpdatePerformanceCycle(PC); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance Cycle Updation failed": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, id)
}

func ListPerformanceCycle(c *gin.Context) {
	PC, err := performance_evaluation.GetPerformanceCycle()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, PC)
}
