// models/employee_document.go
package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeDocument struct {
	DocumentID int       `json:"document_id"`
	EmployeeID int       `json:"employee_id"`
	FilePath   string    `json:"file_path"`
	UploadedAt time.Time `json:"uploaded_at"`
	UploadedBy int       `json:"uploaded_by"`
}

func InsertEmployeeDocument(db *pgxpool.Pool, doc EmployeeDocument) error {
	query := `
		INSERT INTO employee_documents (
			employee_id, file_path, uploaded_at, uploaded_by
		) VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(context.Background(), query,
		doc.EmployeeID, doc.FilePath, doc.UploadedAt, doc.UploadedBy,
	)
	return err
}
