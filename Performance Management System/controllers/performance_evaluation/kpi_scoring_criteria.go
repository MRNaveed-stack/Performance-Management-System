package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddKpiScoringCriteria(c *gin.Context) {
	var KSC performance_evaluation.KpiScoringCriteria

	if err := c.ShouldBindJSON(&KSC); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if err := performance_evaluation.AddKpiScoringCriteria(KSC); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Criteria Insertion failed": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "Kpi Scoring criteria added"})
}

func UpdateKpiScoringCriteria(c *gin.Context) {
	idParam := c.Param("kpi_scoring_criteria_id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var KSC performance_evaluation.KpiScoringCriteria

	if err := c.ShouldBindJSON(&KSC); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	KSC.KpiScoringCriteriaID = id
	if err := performance_evaluation.UpdateKpiScoringCriteria(KSC); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Criteria updation failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Kpi scoring criteria updated"})
}

func ListKpiScoringCriteria(c *gin.Context) {
	KpiID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	criteria, err := performance_evaluation.GetKpiCriteria(KpiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, criteria)
}
