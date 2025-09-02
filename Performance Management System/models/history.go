package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmploymentHistory struct {
	ID           int       `json:"history_id"`
	EmployeeID   int       `json:"employee_id"`
	FieldChanged string    `json:"field_changed"`
	OldValue     string    `json:"old_value"`
	NewValue     string    `json:"new_value"`
	ChangedBy    int       `json:"changed_by"`
	ChangedAt    time.Time `json:"change_timestamp"`
}

type InsertHistoryInput struct {
	EmployeeID   int
	FieldChanged string
	OldValue     string
	NewValue     string
	ChangedBy    int
}

func InsertEmploymentHistory(db *pgxpool.Pool, input InsertHistoryInput) error {
	_, err := db.Exec(context.Background(), `
        INSERT INTO employment_history 
        (employee_id, field_changed, old_value, new_value, changed_by)
        VALUES ($1, $2, $3, $4, $5)
    `, input.EmployeeID, input.FieldChanged, input.OldValue, input.NewValue, input.ChangedBy)
	return err
}

func GetEmployeeHistoryByID(db *pgxpool.Pool, employeeID int) ([]EmploymentHistory, error) {
	rows, err := db.Query(context.Background(),
		`SELECT history_id, employee_id, field_changed, old_value, new_value, changed_by, change_timestamp 
		 FROM employment_history
		 WHERE employee_id = $1 ORDER BY change_timestamp DESC`, employeeID)
	if err != nil {
		log.Printf("Query error for employee_id %d: %v", employeeID, err)
		return nil, err
	}
	defer rows.Close()

	var history []EmploymentHistory
	for rows.Next() {
		var h EmploymentHistory
		err := rows.Scan(&h.ID, &h.EmployeeID, &h.FieldChanged, &h.OldValue, &h.NewValue, &h.ChangedBy, &h.ChangedAt)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		history = append(history, h)
	}

	log.Printf("Found %d history records for employee_id %d", len(history), employeeID)
	return history, nil
}

func GetAllEmploymentHistory(db *pgxpool.Pool) ([]EmploymentHistory, error) {
	rows, err := db.Query(context.Background(),
		`SELECT history_id, employee_id, field_changed, old_value, new_value, changed_by, change_timestamp 
		 FROM employment_history ORDER BY change_timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []EmploymentHistory
	for rows.Next() {
		var h EmploymentHistory
		err := rows.Scan(&h.ID, &h.EmployeeID, &h.FieldChanged, &h.OldValue, &h.NewValue, &h.ChangedBy, &h.ChangedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}
