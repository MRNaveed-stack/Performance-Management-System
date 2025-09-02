package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Education struct {
	EducationID      int    `json:"education_id"`
	EmployeeID       int    `json:"employee_id"`
	DegreeName       string `json:"degree_name"`
	InstituteName    string `json:"institute_name"`
	YearOfCompletion int    `json:"year_of_completion"`
}

// Create new education record
func CreateEducation(db *pgxpool.Pool, edu Education) error {
	query := `
		INSERT INTO educations (employee_id, degree_name, institute_name, year_of_completion)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(context.Background(), query,
		edu.EmployeeID,
		edu.DegreeName,
		edu.InstituteName,
		edu.YearOfCompletion,
	)
	return err
}

// Get all education records for an employee
func GetEducationByEmployeeID(db *pgxpool.Pool, employeeID int) ([]Education, error) {
	query := `
	SELECT education_id, employee_id, degree_name, institute_name, year_of_completion
	FROM educations
	WHERE employee_id = $1
`
	rows, err := db.Query(context.Background(), query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var educations []Education
	for rows.Next() {
		var e Education
		err := rows.Scan(
			&e.EducationID,
			&e.EmployeeID,
			&e.DegreeName,
			&e.InstituteName,
			&e.YearOfCompletion,
		)
		if err != nil {
			return nil, err
		}
		educations = append(educations, e)
	}
	return educations, nil
}

// Delete a specific education record
func DeleteEducationByID(db *pgxpool.Pool, educationID int) error {
	query := `DELETE FROM educations WHERE education_id = $1`
	_, err := db.Exec(context.Background(), query, educationID)
	return err
}
