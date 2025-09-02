package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddPerformanceApproval(c *gin.Context) {
	var PA performance_evaluation.PerformanceApproval

	if err := c.ShouldBindJSON(&PA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.InsertPerformanceApproval(PA); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance Approval Insertion failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Performance Approval Inserted Successfully"})
}

func UPDATEPerformanceApproval(c *gin.Context) {
	var PA performance_evaluation.PerformanceApproval

	idParam := c.Param("approval_id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.ShouldBindJSON(&PA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	PA.ApprovalID = id
	if err := performance_evaluation.UpdatePerformanceApproval(PA); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance Approval Updation failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Performance Approval Updated Successfuly"})
}

func ListPerformanceApproval(c *gin.Context) {
	idParam := c.Param("approval_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	approval, err := performance_evaluation.GetPerformanceApproval(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to fetch the performance approval": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, approval)
}
