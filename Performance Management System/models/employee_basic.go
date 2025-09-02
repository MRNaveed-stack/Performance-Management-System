package models

import (
	"context"
	"performanceManagement/config"
)

type EmployeeBasic struct {
	Name              string `json:"employee_name"`
	FatherHusbandName string `json:"employee_father_husband_name"`
	ContactNumber     string `json:"employee_contact_number"`
	Designation       string `json:"employee_designation"`
}

// This struct was used before the modification in frontend
type EmployeeBasicInfo struct {
	EmployeeID        int    `json:"employee_id"`
	Name              string `json:"employee_name"`
	FatherHusbandName string `json:"employee_father_husband_name"`
	Gender            string `json:"gender"`
	MaritalStatus     string `json:"marital_status"`
	ContactNumber     string `json:"employee_contact_number"`
	Department        string `json:"department"`
	Designation       string `json:"employee_designation"`
}

// This function was used before the modification in frontend
func (e *EmployeeBasic) Create() error {
	query := `
		INSERT INTO employees (
			employee_name,
			employee_father_husband_name,
			employee_contact_number,
			employee_designation
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		e.Name,
		e.FatherHusbandName,
		e.ContactNumber,
		e.Designation,
	)
	return err
}

func GetAllEmployeeBasicInfo() ([]EmployeeBasicInfo, error) {
	query := `
	SELECT 
		employee_id,
		employee_name,
		employee_father_husband_name,
		gender,
		marital_status,
		employee_contact_number,
		department,
		employee_designation
	FROM fn_list_employee_basic()
`

	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []EmployeeBasicInfo

	for rows.Next() {
		var e EmployeeBasicInfo
		if err := rows.Scan(
			&e.EmployeeID,
			&e.Name,
			&e.FatherHusbandName,
			&e.Gender,
			&e.MaritalStatus,
			&e.ContactNumber,
			&e.Department,
			&e.Designation,
		); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

// Before modification, this function was used to update employee basic info
func (e *EmployeeBasic) Update(id string) error {
	query := `
		UPDATE employees
		SET
			employee_name = $1,
			employee_father_husband_name = $2,
			employee_contact_number = $3,
			employee_designation = $4
		WHERE employee_id = $5
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		e.Name,
		e.FatherHusbandName,
		e.ContactNumber,
		e.Designation,
		id,
	)

	return err
}

func InsertEmployeeStoredProc(
	userID int,
	name string,
	fatherHusbandName string,
	contactNumber string,
	genderID int,
	maritalStatusID int,
	departmentID int,
	designation string,
) error {
	query := `CALL sp_insert_employee_basic($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := config.DB.Exec(context.Background(), query,
		userID, name, fatherHusbandName, contactNumber,
		genderID, maritalStatusID, departmentID, designation)
	return err
}

func UpdateEmployeeStoredProc(
	id, name, fatherHusbandName, contactNumber string,
	genderID, maritalStatusID, departmentID int,
	designation string,
) error {
	query := `CALL sp_update_employee_basic($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := config.DB.Exec(context.Background(), query,
		id, name, fatherHusbandName, contactNumber,
		genderID, maritalStatusID, departmentID, designation,
	)
	return err
}

func GetEmployeeBasicInfoByID(employeeID int) (EmployeeBasicInfo, error) {
	var info EmployeeBasicInfo

	query := `SELECT * FROM fn_get_employee_basic_info_by_id($1)`

	err := config.DB.QueryRow(context.Background(), query, employeeID).Scan(
		&info.EmployeeID,
		&info.Name,
		&info.FatherHusbandName,
		&info.Gender,
		&info.MaritalStatus,
		&info.ContactNumber,
		&info.Department,
		&info.Designation,
	)

	return info, err
}
