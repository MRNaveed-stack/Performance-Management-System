package performance_evaluation

import (
	"net/http"
	"path/filepath"
	"performanceManagement/models/performance_evaluation"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddPerformanceDocument(c *gin.Context) {
	employeeID, _ := strconv.Atoi(c.PostForm("employee_id"))
	performanceCycleID, _ := strconv.Atoi(c.PostForm("performance-cycle_id"))
	uploadedBy, _ := strconv.Atoi(c.PostForm("uploaded_by"))
	documentType := c.PostForm("document_type")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"File upload failed": err.Error()})
		return
	}

	uploadPath := "./uploads" + file.Filename
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error while saving file": err.Error()})
		return
	}
	fileURL := "/uploads/" + file.Filename

	PD := performance_evaluation.PerformanceDocument{
		EmployeeID:         employeeID,
		PerformanceCycleID: performanceCycleID,
		UploadedBy:         uploadedBy,
		DocumentType:       documentType,
		FileURL:            fileURL,
	}

	if err := performance_evaluation.InsertPerformanceDocument(PD); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Document Uploaded Successfully"})
}

func ListPerformanceDocuments(c *gin.Context) {
	empidParam := c.Query("employee_id")
	performanceidParam := c.Query("performance_cycle_id")

	var empID, cycleID *int
	if empidParam != "" {
		val, _ := strconv.Atoi(empidParam)
		empID = &val
	}
	if performanceidParam != "" {
		val, _ := strconv.Atoi(performanceidParam)
		cycleID = &val
	}

	docs, err := performance_evaluation.GetPerformanceDocument(*empID, *cycleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, docs)
}

func DownloadPerformanceDocument(c *gin.Context) {
	idParam := c.Param("document_id")
	id, _ := strconv.Atoi(idParam)
	doc, err := performance_evaluation.GetPerformanceDocumentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment("."+doc.FileURL, filepath.Base(doc.FileURL))
}

func DeletePerformanceDocument(c *gin.Context) {
	idParam := c.Param("document_id")
	id, _ := strconv.Atoi(idParam)
	if err := performance_evaluation.RemovePerformanceDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Message": "Performance Document Deleted Successfully"})
}
