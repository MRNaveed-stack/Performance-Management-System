package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddSelfAppraisal(c *gin.Context) {
	var SA performance_evaluation.SelfAppraisal

	if err := c.ShouldBindJSON(&SA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.InsertSelfAppraisal(SA); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Self Appraisal Inertion failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Self Appraisal Inserted Successfully"})
}

func UPDATESelfAppraisal(c *gin.Context) {
	idParam := c.Param("self_appraisal_id")
	var SA performance_evaluation.SelfAppraisal

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.ShouldBindJSON(&SA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	SA.SelfAppraisalID = id

	if err := performance_evaluation.UpdateSelfAppraisal(SA); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Self Appraisal Updation failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Self appraisal updated successfully"})
}

func ListSelfAppraisal(c *gin.Context) {
	appraisals, err := performance_evaluation.GetPerformanceCycle()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to fetch self appraisals": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, appraisals)
}

func ListSelfAppraisalByID(c *gin.Context) {
	idParam := c.Param("self_appraisal_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	appraisal, err := performance_evaluation.GetSelfAppraisalByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to fetch self appraisal": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, appraisal)
}
