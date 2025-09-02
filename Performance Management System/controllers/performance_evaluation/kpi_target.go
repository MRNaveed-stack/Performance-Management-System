package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddKpiTarget(c *gin.Context) {
	var KPIt performance_evaluation.KpiTarget

	if err := c.ShouldBindJSON(&KPIt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.InsertKpiTarget(KPIt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "Kpi target value added successfully"})
}

func UPDATEKpiTarget(c *gin.Context) {
	idParam := c.Param("kpi_target_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalidID"})
		return
	}

	var K performance_evaluation.KpiTarget

	if err := c.ShouldBindJSON(&K); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	K.KpiTargetID = id

	if err := performance_evaluation.UpdateKpiTarget(K); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Kpi Target Updation Failed": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Kpi target updated successfully"})
}

func ListKpiTarget(c *gin.Context) {
	targets, err := performance_evaluation.GetKpiTarget()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to fetch the kpi target": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, targets)
}
