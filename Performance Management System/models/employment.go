package models

import (
	"context"
	"errors"
	"performanceManagement/config"
	"time"
)

type EmploymentDetails1 struct {
	EmployeeID         int       `json:"employee_id" binding:"required"`
	JoiningDate        time.Time `json:"employee_joining_date" binding:"required"`
	EmploymentTypeID   int       `json:"employment_type_id" binding:"required"`
	EmploymentStatusID int       `json:"employment_status_id" binding:"required"`
}

type EmploymentDetails struct {
	JoiningDate      string `json:"employee_joining_date"`
	EmploymentType   string `json:"employment_type"`
	EmploymentStatus string `json:"employment_status"`
}

type EmploymentDetailsUpdate struct {
	JoiningDate        *string `json:"employee_joining_date,omitempty"`
	EmploymentTypeID   *int    `json:"employment_type_id,omitempty"`
	EmploymentStatusID *int    `json:"employment_status_id,omitempty"`
}

func GetEmploymentDetailsByID(empID int) (*EmploymentDetails, error) {
	query := `
	SELECT 
		TO_CHAR(employee_joining_date, 'YYYY-MM-DD') AS employee_joining_date,
		et.employment_type_name,
		es.employment_status_name
	FROM 
		employees e
	JOIN 
		employment_type et ON e.employment_type_id = et.employment_type_id
	JOIN 
		employment_status es ON e.employee_status_id = es.employment_status_id
	WHERE 
		e.employee_id = $1
	`

	var details EmploymentDetails
	err := config.DB.QueryRow(context.Background(), query, empID).Scan(
		&details.JoiningDate,
		&details.EmploymentType,
		&details.EmploymentStatus,
	)

	if err != nil {
		return nil, err
	}

	return &details, nil
}

func UpdateEmploymentDetails(empID int, input EmploymentDetailsUpdate) error {
	query := `
UPDATE employees
SET 
    employee_joining_date = COALESCE($1, employee_joining_date),
    employment_type_id = COALESCE($2, employment_type_id),
    employee_status_id = COALESCE($3, employee_status_id)
WHERE employee_id = $4
`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		input.JoiningDate,
		input.EmploymentTypeID,
		input.EmploymentStatusID,
		empID,
	)

	return err
}

func (e *EmploymentDetails1) SaveEmploymentDetails() error {
	query := `
		UPDATE employees
		SET 
			employee_joining_date = $1,
			employment_type_id = $2,
			employee_status_id = $3
		WHERE employee_id = $4
	`

	cmdTag, err := config.DB.Exec(
		context.Background(),
		query,
		e.JoiningDate,
		e.EmploymentTypeID,
		e.EmploymentStatusID,
		e.EmployeeID,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no employee found with the given ID")
	}

	return nil
}

type EmploymentType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EmploymentStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllEmploymentTypes() ([]EmploymentType, error) {
	rows, err := config.DB.Query(context.Background(), `SELECT employment_type_id, employment_type_name FROM employment_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []EmploymentType
	for rows.Next() {
		var t EmploymentType
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func GetAllEmploymentStatuses() ([]EmploymentStatus, error) {
	rows, err := config.DB.Query(context.Background(), `SELECT employment_status_id, employment_status_name FROM employment_status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []EmploymentStatus
	for rows.Next() {
		var s EmploymentStatus
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}
