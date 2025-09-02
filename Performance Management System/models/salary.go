package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SalaryInfo struct {
	SalaryID      int       `json:"salary_id"`
	EmployeeID    int       `json:"employee_id"`
	PayGrade      string    `json:"pay_grade"`
	SalaryBand    string    `json:"salary_band"` // <- updated
	SalaryAmount  float64   `json:"salary_amount"`
	Reason        string    `json:"reason"`
	UpdatedBy     int       `json:"updated_by"`
	Actions       string    `json:"actions"`
	EffectiveDate time.Time `json:"effective_date"`
	ActionID      int       `json:"action_id"`
}

func CreateSalaryInfo(db *pgxpool.Pool, salary SalaryInfo) error {
	query := `
	INSERT INTO salary_info (
		employee_id, pay_grade, salary_band,
		salary_amount, reason, updated_by, actions, effective_date, action_id
	) VALUES (
		$1, $2, $3,
		$4, $5, $6, $7, $8, $9
	)
`

	_, err := db.Exec(
		context.Background(),
		query,
		salary.EmployeeID,
		salary.PayGrade,
		salary.SalaryBand, // <- updated
		salary.SalaryAmount,
		salary.Reason,
		salary.UpdatedBy,
		salary.Actions,
		salary.EffectiveDate,
		salary.ActionID,
	)

	return err
}
func GetSalaryInfoByEmployeeID(db *pgxpool.Pool, employeeID int) ([]SalaryInfo, error) {
	query := `
	SELECT salary_id, employee_id, pay_grade, salary_band,
	       salary_amount, reason, updated_by, actions, effective_date, action_id
	FROM salary_info
	WHERE employee_id = $1
	ORDER BY effective_date DESC
`

	rows, err := db.Query(context.Background(), query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salaries []SalaryInfo
	for rows.Next() {
		var s SalaryInfo
		err := rows.Scan(
			&s.SalaryID,
			&s.EmployeeID,
			&s.PayGrade,
			&s.SalaryBand,

			&s.SalaryAmount,
			&s.Reason,
			&s.UpdatedBy,
			&s.Actions,
			&s.EffectiveDate,
			&s.ActionID,
		)
		if err != nil {
			return nil, err
		}
		salaries = append(salaries, s)
	}

	return salaries, nil
}

func UpdateSalaryInfo(db *pgxpool.Pool, salary SalaryInfo) error {
	query := `
	UPDATE salary_info
	SET pay_grade = $1,
	    salary_band = $2,
	    salary_amount = $3,
	    reason = $4,
	    updated_by = $5,
	    actions = $6,
	    effective_date = $7,
	    action_id = $8
	WHERE salary_id = $9
`

	_, err := db.Exec(
		context.Background(),
		query,
		salary.PayGrade,
		salary.SalaryBand, // <- updated
		salary.SalaryAmount,
		salary.Reason,
		salary.UpdatedBy,
		salary.Actions,
		salary.EffectiveDate,
		salary.ActionID,
		salary.SalaryID,
	)

	return err
}

func DeleteSalaryInfo(db *pgxpool.Pool, salaryID int) error {
	query := `DELETE FROM salary_info WHERE salary_id = $1`

	_, err := db.Exec(context.Background(), query, salaryID)
	return err
}
