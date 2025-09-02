package performance_evaluation

import (
	"net/http"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddKpiCategory(c *gin.Context) {

	var category performance_evaluation.KpiCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := performance_evaluation.InsertKpiCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Insertion failed": err.Error()})
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Kpi category added successfully"})
}

func UpdateKpiCategory(c *gin.Context) {
	idParam := c.Param("kpi_category_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	var category performance_evaluation.KpiCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	category.CategoryID = id
	if err := performance_evaluation.UpdateKpiCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Updation failed": err.Error()})
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Kpi category updated successfully"})
}

func GetKpiCategory(c *gin.Context) {
	categories, err := performance_evaluation.ListKpiCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusAccepted, categories)

}
