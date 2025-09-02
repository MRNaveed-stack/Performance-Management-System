package controllers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"performanceManagement/config"
	"performanceManagement/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"log"
)

func CreateEmployeeBasic(c *gin.Context) {
	// This struct combines basic info + dropdowns
	var input struct {
		Name              string `json:"employee_name"`
		FatherHusbandName string `json:"employee_father_husband_name"`
		ContactNumber     string `json:"employee_contact_number"`
		GenderID          int    `json:"gender_id"`
		MaritalStatusID   int    `json:"marital_status_id"`
		DepartmentID      int    `json:"employee_department_id"`
		Designation       string `json:"employee_designation"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	var userID int
	switch v := userIDInterface.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id invalid type"})
		return
	}

	err := models.InsertEmployeeStoredProc(
		userID,
		input.Name,
		input.FatherHusbandName,
		input.ContactNumber,
		input.GenderID,
		input.MaritalStatusID,
		input.DepartmentID,
		input.Designation,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert employee", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Employee created successfully"})
}

func GetAllEmployeeBasicInfo(c *gin.Context) {
	employees, err := models.GetAllEmployeeBasicInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch employees"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func GetAllGenders(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT gender_id, gender_name FROM genders")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genders"})
		return
	}
	defer rows.Close()

	var genders []gin.H
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		genders = append(genders, gin.H{"id": id, "name": name})
	}

	c.JSON(http.StatusOK, genders)
}

func GetAllMaritalStatuses(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT marital_status_id, marital_status_name FROM marital_statuses")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch marital statuses"})
		return
	}
	defer rows.Close()

	var statuses []gin.H
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		statuses = append(statuses, gin.H{"id": id, "name": name})
	}

	c.JSON(http.StatusOK, statuses)
}

func GetAllDepartments(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT department_id, department_name FROM departments")
	if err != nil {
		fmt.Println("DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}
	defer rows.Close()

	var departments []gin.H
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		departments = append(departments, gin.H{"id": id, "name": name})
	}

	c.JSON(http.StatusOK, departments)
}

func GetEmployeeBasicInfoByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var info models.EmployeeBasicInfo

	query := `SELECT * FROM fn_get_employee_basic_info_by_id($1)`

	err = config.DB.QueryRow(context.Background(), query, id).Scan(
		&info.EmployeeID,
		&info.Name,
		&info.FatherHusbandName,
		&info.Gender,
		&info.MaritalStatus,
		&info.ContactNumber,
		&info.Department,
		&info.Designation,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

func UpdateEmployeeBasicInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID is required"})
		return
	}

	var input struct {
		Name              string `json:"employee_name"`
		FatherHusbandName string `json:"employee_father_husband_name"`
		ContactNumber     string `json:"employee_contact_number"`
		GenderID          int    `json:"gender_id"`
		MaritalStatusID   int    `json:"marital_status_id"`
		DepartmentID      int    `json:"employee_department_id"`
		Designation       string `json:"employee_designation"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	err := models.UpdateEmployeeStoredProc(
		id,
		input.Name,
		input.FatherHusbandName,
		input.ContactNumber,
		input.GenderID,
		input.MaritalStatusID,
		input.DepartmentID,
		input.Designation,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

func GetEmploymentDetails(c *gin.Context) {
	idParam := c.Param("id")
	empID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	details, err := models.GetEmploymentDetailsByID(empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch employment details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": details})
}

func UpdateEmploymentDetails(c *gin.Context) {
	idParam := c.Param("id")
	empID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	var input models.EmploymentDetailsUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	if err := models.UpdateEmploymentDetails(empID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update employment details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employment details updated successfully"})
}

func SaveEmploymentDetailsHandler(c *gin.Context) {
	idParam := c.Param("id")
	employeeID, err := strconv.Atoi(idParam)
	if err != nil || employeeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req struct {
		JoiningDate        string `json:"employee_joining_date" binding:"required"`
		EmploymentTypeID   int    `json:"employment_type_id" binding:"required"`
		EmploymentStatusID int    `json:"employee_status_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	joiningDate, err := time.Parse("2006-01-02", req.JoiningDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid joining date format. Use YYYY-MM-DD"})
		return
	}

	employment := models.EmploymentDetails1{
		EmployeeID:         employeeID,
		JoiningDate:        joiningDate,
		EmploymentTypeID:   req.EmploymentTypeID,
		EmploymentStatusID: req.EmploymentStatusID,
	}

	if err := employment.SaveEmploymentDetails(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save employment details", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employment details saved successfully"})
}

func GetAllEmploymentTypes(c *gin.Context) {
	data, err := models.GetAllEmploymentTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employment types"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func GetAllEmploymentStatuses(c *gin.Context) {
	data, err := models.GetAllEmploymentStatuses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employment statuses"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// REPORTING STRUCTURE
func GetOwnReportingStructure(c *gin.Context) {
	empIDRaw, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	empID := empIDRaw.(int)

	rs, err := models.GetReportingStructureByID(empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reporting structure"})
		return
	}

	c.JSON(http.StatusOK, rs)
}

func GetAllTeams(c *gin.Context) {
	teams, err := models.FetchAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}
	c.JSON(http.StatusOK, teams)
}

func GetReportingStructureByID(c *gin.Context) {
	idParam := c.Param("id")
	empID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	rs, err := models.GetReportingStructureByID(empID)
	if err != nil {
		fmt.Println("DEBUG: error from GetReportingStructureByID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reporting structure"})
		return
	}

	c.JSON(http.StatusOK, rs)
}
func GetReportingManagers(c *gin.Context) {
	managers, err := models.FetchReportingManagers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch reporting managers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, managers)
}
func GetFullReportingStructureByID(c *gin.Context) {
	idParam := c.Param("id")
	employeeID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Get the main reporting structure
	structure, err := models.GetReportingStructureByID(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reporting structure"})
		return
	}

	// Get all managers for dropdown
	managers, err := models.FetchReportingManagers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch managers"})
		return
	}

	// Get all teams for dropdown
	teams, err := models.FetchAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	// Combine all into one response
	c.JSON(http.StatusOK, gin.H{
		"reporting_structure": structure,
		"managers":            managers,
		"teams":               teams,
	})
}

type UpdateReportingRequest struct {
	ReportingManagerName string `json:"reporting_manager_name"`
	TeamName             string `json:"team_name"`
}

func UpdateReportingStructure(c *gin.Context) {
	idParam := c.Param("id")
	empID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	var input UpdateReportingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := models.SetReportingStructureByNames(empID, input.ReportingManagerName, input.TeamName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reporting structure updated successfully"})
}

// SALARY INFORMATION

type SalaryInfoInput struct {
	EmployeeID    int     `json:"employee_id" binding:"required"`
	PayGrade      string  `json:"pay_grade" binding:"required"`
	SalaryBand    string  `json:"salary_band" binding:"required"` // CHANGED
	SalaryAmount  float64 `json:"salary_amount" binding:"required"`
	Reason        string  `json:"reason"`
	UpdatedBy     int     `json:"updated_by" binding:"required"`
	ActionID      int     `json:"action_id" binding:"required"`
	EffectiveDate string  `json:"effective_date"`
}

func CreateSalaryInfo(c *gin.Context) {
	var input SalaryInfoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	effectiveDate := time.Now()
	if input.EffectiveDate != "" {
		var err error
		effectiveDate, err = time.Parse("2006-01-02", input.EffectiveDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid effective_date format. Use YYYY-MM-DD"})
			return
		}
	}

	// Map input to model
	salary := models.SalaryInfo{
		EmployeeID: input.EmployeeID,
		PayGrade:   input.PayGrade,
		SalaryBand: input.SalaryBand,

		SalaryAmount:  input.SalaryAmount,
		Reason:        input.Reason,
		UpdatedBy:     input.UpdatedBy,
		Actions:       "View", // Default action
		EffectiveDate: effectiveDate,
		ActionID:      input.ActionID,
	}

	// Insert into database using model function
	if err := models.CreateSalaryInfo(config.DB, salary); err != nil {
		fmt.Println("Insert Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Salary info added successfully"})
}

func GetSalaryInfoByEmployeeID(c *gin.Context) {
	employeeIDParam := c.Param("id")
	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	salaryInfoList, err := models.GetSalaryInfoByEmployeeID(config.DB, employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch salary info"})
		return
	}

	c.JSON(http.StatusOK, salaryInfoList)
}

func UpdateSalaryInfo(c *gin.Context) {
	salaryIDParam := c.Param("salary_id")
	salaryID, err := strconv.Atoi(salaryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid salary ID"})
		return
	}

	var input models.SalaryInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.SalaryID = salaryID
	if err := models.UpdateSalaryInfo(config.DB, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update salary info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Salary info updated successfully"})
}

func DeleteSalaryInfo(c *gin.Context) {
	salaryIDParam := c.Param("salary_id")
	salaryID, err := strconv.Atoi(salaryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid salary ID"})
		return
	}

	if err := models.DeleteSalaryInfo(config.DB, salaryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete salary info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Salary info deleted successfully"})
}

// Add educations
func AddEducation(c *gin.Context) {
	var edu models.Education

	// Get employee ID from URL and assign
	employeeIDStr := c.Param("id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := c.ShouldBindJSON(&edu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	edu.EmployeeID = employeeID

	err = models.CreateEducation(config.DB, edu)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Education added successfully"})
}
func GetEducation(c *gin.Context) {
	employeeIDStr := c.Param("id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	educations, err := models.GetEducationByEmployeeID(config.DB, employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch education"})
		return
	}

	c.JSON(http.StatusOK, educations)
}

func DeleteEducation(c *gin.Context) {
	educationIDStr := c.Param("education_id")
	educationID, err := strconv.Atoi(educationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid education ID"})
		return
	}

	err = models.DeleteEducationByID(config.DB, educationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete education"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Education deleted successfully"})
}

// Certifications
func AddCertification(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	var cert models.Certification
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cert.EmployeeID = employeeID

	err = models.CreateCertification(cert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create certification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "certification added"})
}

func GetCertifications(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	certifications, err := models.GetCertificationsByEmployee(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certifications"})
		return
	}

	c.JSON(http.StatusOK, certifications)
}

func DeleteCertification(c *gin.Context) {
	certID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certification id"})
		return
	}

	err = models.DeleteCertification(certID)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "certification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification deleted"})
}

// SKILLS

// POST /skills - Add a new skill (if not already in db)
func CreateSkillHandler(c *gin.Context) {
	var input struct {
		SkillName string `json:"skill_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Skill name is required"})
		return
	}

	id, err := models.CreateSkill(config.DB, input.SkillName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill added", "skill_id": id})
}

// POST /employees/:id/skills - Assign skill to employee
func AddEmployeeSkillHandler(c *gin.Context) {
	empID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var input struct {
		SkillID int `json:"skill_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Skill ID required"})
		return
	}

	err = models.AddEmployeeSkill(config.DB, empID, input.SkillID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill assigned"})
}

// GET /employees/:id/skills - Get all skills for employee
func GetEmployeeSkillsHandler(c *gin.Context) {
	empID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	skills, err := models.GetSkillsByEmployee(config.DB, empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}

	c.JSON(http.StatusOK, skills)
}

// DELETE /employees/:id/skills/:skill_id - Remove specific skill from employee
func DeleteEmployeeSkillHandler(c *gin.Context) {
	empID, err1 := strconv.Atoi(c.Param("id"))
	skillID, err2 := strconv.Atoi(c.Param("skill_id"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee or skill ID"})
		return
	}

	err := models.DeleteEmployeeSkill(config.DB, empID, skillID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill removed"})
}

// Employee documents
func UploadEmployeeDocument(c *gin.Context) {
	// Parse employee_id from URL
	employeeIDStr := c.Param("id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Parse uploaded_by
	uploadedByStr := c.PostForm("uploaded_by")
	uploadedBy, err := strconv.Atoi(uploadedByStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uploaded_by"})
		return
	}

	// Get the file
	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	// Create uploads directory
	uploadPath := "uploads"
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Save the file
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
	filePath := filepath.Join(uploadPath, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	src, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file stream"})
		return
	}
	defer src.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
		return
	}

	// Save metadata in DB (no documentTypeID anymore)
	doc := models.EmployeeDocument{
		EmployeeID: employeeID,
		FilePath:   fileName, // Save relative path only
		UploadedAt: time.Now(),
		UploadedBy: uploadedBy,
	}

	err = models.InsertEmployeeDocument(config.DB, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document uploaded successfully", "file_path": fileName})
}

type EmployeeDocumentResponse struct {
	DocumentID int    `json:"document_id"`
	FilePath   string `json:"file_path"`
	UploadedAt string `json:"uploaded_at"`
	UploadedBy string `json:"uploaded_by"`
}

func GetEmployeeDocuments(c *gin.Context) {
	employeeID := c.Param("id")

	rows, err := config.DB.Query(c, `
		SELECT 
			ed.document_id,
			ed.file_path,
			ed.uploaded_at,
			e.name AS uploaded_by
		FROM employee_documents ed
		JOIN employees e ON ed.uploaded_by = e.employee_id
		WHERE ed.employee_id = $1
	`, employeeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}
	defer rows.Close()

	var documents []EmployeeDocumentResponse

	for rows.Next() {
		var doc EmployeeDocumentResponse
		err := rows.Scan(&doc.DocumentID, &doc.FilePath, &doc.UploadedAt, &doc.UploadedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// GetSingleEmployeeDocument serves the file based on document ID
func GetSingleEmployeeDocument(c *gin.Context) {
	docID := c.Param("doc_id")

	var filePath string
	err := config.DB.QueryRow(c, `
		SELECT file_path FROM employee_documents WHERE document_id = $1
	`, docID).Scan(&filePath)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Construct full path
	fullPath := filepath.Join("uploads", filePath)

	// Optional: Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on server"})
		return
	}

	// Serve the file
	c.File(fullPath)
}

func DeleteEmployeeDocument(c *gin.Context) {
	docID := c.Param("doc_id")

	var filePath string
	err := config.DB.QueryRow(c, `
		SELECT file_path FROM employee_documents WHERE document_id = $1
	`, docID).Scan(&filePath)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	fullPath := filepath.Join("uploads", filePath)

	// Delete file from filesystem
	err = os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from server"})
		return
	}

	// Delete record from DB
	_, err = config.DB.Exec(c, `
		DELETE FROM employee_documents WHERE document_id = $1
	`, docID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
func GetEmployeeHistoryByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		log.Printf("Fetching history for employee ID: %d", id)

		return
	}

	history, err := models.GetEmployeeHistoryByID(config.DB, id)
	if err != nil {
		log.Println("Employment history fetch error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
func InsertEmploymentHistoryHandler(c *gin.Context) {
	var input struct {
		EmployeeID   int    `json:"employee_id"`
		FieldChanged string `json:"field_changed"`
		OldValue     string `json:"old_value"`
		NewValue     string `json:"new_value"`
		ChangedBy    int    `json:"changed_by"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := models.InsertEmploymentHistory(config.DB, models.InsertHistoryInput{
		EmployeeID:   input.EmployeeID,
		FieldChanged: input.FieldChanged,
		OldValue:     input.OldValue,
		NewValue:     input.NewValue,
		ChangedBy:    input.ChangedBy,
	})

	if err != nil {
		log.Println("InsertEmploymentHistory error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Employment history recorded"})
}

func GetAllEmploymentHistoryHandler(c *gin.Context) {
	history, err := models.GetAllEmploymentHistory(config.DB)

	if err != nil {
		fmt.Println("ERROR fetching history:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
