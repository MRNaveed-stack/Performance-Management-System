package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddPerformanceReview(c *gin.Context) {
	var PR performance_evaluation.PerformanceReview

	if err := c.ShouldBindJSON(&PR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.InsertPerformanceReview(PR); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance RevieW Insertion Failed": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "Performance Review Inserted Successfully"})
}

func UPDATEPerformanceReview(c *gin.Context) {
	var PR performance_evaluation.PerformanceReview
	idParam := c.Param("performance_review_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.ShouldBindJSON(&PR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	PR.PerformanceReviewID = id

	if err := performance_evaluation.UpdatePerformanceReview(PR); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Performance Review Updation Failed": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "Performance review updated successfully"})
}

func ListPerformanceReview(c *gin.Context) {
	review, err := performance_evaluation.GetPerformanceReview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, review)
}

func ListPerformanceReviewByID(c *gin.Context) {
	idParam := c.Param("performance_review_id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	review, err := performance_evaluation.GetPerformanceReviewByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to fetch performance review": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, review)
}
