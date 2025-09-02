package attendance

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"performanceManagement/config"
	"performanceManagement/models/attendance"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InsertLeaveApplication(c *gin.Context) {
	// Parse form fields (not JSON!)
	employeeID, _ := strconv.Atoi(c.PostForm("employee_id"))
	leaveTypeID, _ := strconv.Atoi(c.PostForm("leave_type_id"))
	departmentID, _ := strconv.Atoi(c.PostForm("department_id"))
	currentStatusID, _ := strconv.Atoi(c.PostForm("current_status_id"))

	title := c.PostForm("title")
	description := c.PostForm("description")
	fromDate := c.PostForm("from_date")
	toDate := c.PostForm("to_date")

	var attachedDocument string
	file, err := c.FormFile("attached_document")
	if err == nil {

		filePath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		attachedDocument = filePath
	} else {
		attachedDocument = ""
	}

	_, err = config.DB.Exec(
		c,
		`CALL sp_insert_leave_application($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		employeeID,
		leaveTypeID,
		departmentID,
		title,
		description,
		fromDate,
		toDate,
		attachedDocument,
		currentStatusID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Leave application submitted successfully"})
}
func GetLeaveApplicationByID(c *gin.Context) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave application ID"})
		return
	}

	var la attendance.LeaveApplication
	query := `
		SELECT leave_application_id, employee_id, leave_type_id, department_id,
		       title, description, from_date, to_date, number_of_leaves,
		       attached_document, created_date, current_status_id
		FROM leave_applications
		WHERE leave_application_id = $1
	`

	err = config.DB.QueryRow(c, query, id).Scan(
		&la.LeaveApplicationID,
		&la.EmployeeID,
		&la.LeaveTypeID,
		&la.DepartmentID,
		&la.Title,
		&la.Description,
		&la.FromDate,
		&la.ToDate,
		&la.NumberOfLeaves,
		&la.AttachedDocument,
		&la.CreatedDate,
		&la.CurrentStatusID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No leave application found for ID %d", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave application not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	if la.CurrentStatusID == 0 {
		log.Printf("Leave application %d has no status assigned", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave application has no status assigned"})
		return
	}

	log.Printf("Successfully fetched leave application ID %d", id)
	c.JSON(http.StatusOK, la)
}

func UpdateLeaveApplicationByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave application ID"})
		return
	}

	// Parse form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	employeeID, _ := strconv.Atoi(c.PostForm("employee_id"))
	departmentID, _ := strconv.Atoi(c.PostForm("department_id"))
	leaveTypeID, _ := strconv.Atoi(c.PostForm("leave_type_id"))
	title := c.PostForm("title")
	description := c.PostForm("description")
	fromDate := c.PostForm("from_date")
	toDate := c.PostForm("to_date")
	currentStatusID, _ := strconv.Atoi(c.PostForm("current_status_id"))

	var filePath string
	file, err := c.FormFile("attached_document")
	if err == nil {
		filePath = "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
			return
		}
	} else {
		filePath = c.PostForm("attached_document")
	}

	query := `
		CALL sp_update_leave_application(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	_, err = config.DB.Exec(
		c, query,
		id,
		employeeID,
		departmentID,
		leaveTypeID,
		title,
		description,
		fromDate,
		toDate,
		filePath,
		currentStatusID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave application updated successfully"})
}

func ApplyLeave(c *gin.Context) {
	employeeID, _ := strconv.Atoi(c.PostForm("employee_id"))
	leaveTypeID, _ := strconv.Atoi(c.PostForm("leave_type_id"))
	departmentID, _ := strconv.Atoi(c.PostForm("department_id"))

	const pendingManagerID = 1

	title := c.PostForm("title")
	description := c.PostForm("description")
	fromDate := c.PostForm("from_date")
	toDate := c.PostForm("to_date")

	var attachedDocument string
	file, err := c.FormFile("attached_document")
	if err == nil {
		filePath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		attachedDocument = filePath
	} else {
		attachedDocument = ""
	}

	// Call your stored procedure
	_, err = config.DB.Exec(
		c,
		`CALL sp_insert_leave_application($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		employeeID,
		leaveTypeID,
		departmentID,
		title,
		description,
		fromDate,
		toDate,
		attachedDocument,
		pendingManagerID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Values: %+v", err)

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Leave application submitted and pending manager approval"})
}

func GetLeavesForManager(c *gin.Context) {
	managerIDStr := c.Param("manager_id")
	managerID, err := strconv.Atoi(managerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID"})
		return
	}

	query := `
	SELECT 
		la.leave_application_id,
		la.employee_id,
		la.leave_type_id,
		e.employee_department_id,
		la.title,
		la.description,
		la.from_date,
		la.to_date,
		la.attached_document,
		la.current_status_id,
		la.reporting_manager_status,
		la.reporting_manager_remarks,
		la.hr_status,
		la.hr_remarks
	FROM leave_applications la
	JOIN employees e ON la.employee_id = e.employee_id
	WHERE e.reporting_manager_id = $1 
	  AND la.reporting_manager_status = 'Pending'
	`

	rows, err := config.DB.Query(context.Background(), query, managerID)
	if err != nil {
		log.Printf("DB query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query leaves"})
		return
	}
	defer rows.Close()

	var leaves []attendance.LeaveApplicationRequest

	for rows.Next() {
		var leave attendance.LeaveApplicationRequest

		err := rows.Scan(
			&leave.LeaveApplicationID,
			&leave.EmployeeID,
			&leave.LeaveTypeID,
			&leave.DepartmentID,
			&leave.Title,
			&leave.Description,
			&leave.FromDate,
			&leave.ToDate,
			&leave.AttachedDocument,
			&leave.CurrentStatusID,
			&leave.ReportingManagerStatus,
			&leave.ReportingManagerRemarks,
			&leave.HRStatus,
			&leave.HRRemarks,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		leaves = append(leaves, leave)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Iteration error"})
		return
	}

	c.JSON(http.StatusOK, leaves)
}

func ManagerApproveLeave(c *gin.Context) {
	leaveID := c.Param("leave_id")
	var input struct {
		Status  string `json:"status"` // "Approved" or "Rejected"
		Remarks string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `
	UPDATE leave_applications
	SET reporting_manager_status = $1,
		reporting_manager_remarks = $2
	WHERE leave_application_id = $3
	`

	_, err := config.DB.Exec(context.Background(), query, input.Status, input.Remarks, leaveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update manager decision"})
		return
	}

	if input.Status != "Approved" && input.Status != "Rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager action recorded"})
}

func GetLeavesForHR(c *gin.Context) {
	query := `
	SELECT * FROM leave_applications
	WHERE reporting_manager_status = 'Approved'
	  AND hr_status = 'Pending'
	`

	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	var leaves []attendance.LeaveApplicationRequest1
	for rows.Next() {
		var leave attendance.LeaveApplicationRequest1
		if err := rows.Scan(
			&leave.LeaveApplicationID,
			&leave.EmployeeID,
			&leave.LeaveTypeID,
			&leave.DepartmentID,
			&leave.Title,
			&leave.Description,
			&leave.FromDate,
			&leave.ToDate,
			&leave.NumberOfLeaves, // Add this field
			&leave.AttachedDocument,
			&leave.CreatedDate, // Add this field
			&leave.CurrentStatusID,
			&leave.ReportingManagerStatus,
			&leave.ReportingManagerRemarks,
			&leave.HRStatus,
			&leave.HRRemarks,
		); err != nil {
			log.Printf("Scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		leaves = append(leaves, leave)
	}

	c.JSON(http.StatusOK, leaves)
}

func HRApproveLeave(c *gin.Context) {
	leaveID := c.Param("leave_id")
	var input struct {
		Status  string `json:"status"`
		Remarks string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `
	UPDATE leave_applications
	SET hr_status = $1,
		hr_remarks = $2
	WHERE leave_application_id = $3
	`

	_, err := config.DB.Exec(context.Background(), query, input.Status, input.Remarks, leaveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update HR decision"})
		return
	}

	if input.Status != "Approved" && input.Status != "Rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "HR action recorded"})
}
