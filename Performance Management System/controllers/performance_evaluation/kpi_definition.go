package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddKpi(c *gin.Context) {
	var kpi performance_evaluation.Kpi

	if err := c.ShouldBindJSON(&kpi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := performance_evaluation.AddKpi(kpi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Kpi insertion failed": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "Kpi added successfully"})
}

func UpdateKpi(c *gin.Context) {
	idParam := c.Param("kpi_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var kpi performance_evaluation.Kpi
	if err := c.ShouldBindJSON(&kpi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	kpi.KpiID = id
	if err := performance_evaluation.UpdateKpi(kpi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"kpi updation failed": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "Kpi updated successfully"})
}

func ListKpi(c *gin.Context) {
	var KpiCategoryID *int

	if catID := c.Query("kpi_category_id"); catID != "" {
		if id, err := strconv.Atoi(catID); err == nil {
			KpiCategoryID = &id
		}

	}
	kpi, err := performance_evaluation.ListKpi(KpiCategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, kpi)
}
