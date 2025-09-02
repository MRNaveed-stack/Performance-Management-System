package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type PerformanceDocument struct {
	DocumentID         int       `json:"document_id"`
	EmployeeID         int       `json:"employee_id"`
	PerformanceCycleID int       `json:"performance_cycle_id"`
	UploadedBy         int       `json:"uploaded_by"`
	DocumentType       string    `json:"document_type"`
	FileURL            string    `json:"file_url"`
	UploadedAt         time.Time `json:"uploaded_at"`
}

func InsertPerformanceDocument(PD PerformanceDocument) error {
	_, err := config.DB.Exec(context.Background(),
		`SELECT sp_insert_performance_document($1,$2,$3,$4,$5)`,
		PD.EmployeeID, PD.PerformanceCycleID, PD.UploadedBy, PD.DocumentType, PD.FileURL)
	return err
}

func UpdatePerformanceDocument(PD PerformanceDocument) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_update_performance_document($1,$2,$3)`,
		PD.DocumentID, PD.DocumentType, PD.FileURL)
	return err
}

func GetPerformanceDocument(PerformanceCycleID, EmployeeID int) ([]PerformanceDocument, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * from fn_list_performance_documents($1,$2)`, PerformanceCycleID, EmployeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performancedocument []PerformanceDocument
	for rows.Next() {
		var PD PerformanceDocument
		err := rows.Scan(&PD.DocumentID, &PD.EmployeeID, &PD.PerformanceCycleID, &PD.UploadedBy,
			&PD.DocumentType, &PD.FileURL, &PD.UploadedAt)
		if err != nil {
			return nil, err
		}
		performancedocument = append(performancedocument, PD)
	}
	return performancedocument, nil
}

func GetPerformanceDocumentByID(documentID int) (PerformanceDocument, error) {
	var PD PerformanceDocument

	err := config.DB.QueryRow(context.Background(),
		`SELECT document_id,employee_id,performance_cycle_id,uploaded_by,document_type,file_url,uploaded_at
	 FROM performance_documents where document_id = $1`, documentID).Scan(&PD.DocumentID,
		&PD.EmployeeID, &PD.PerformanceCycleID, &PD.UploadedBy, &PD.DocumentType, &PD.FileURL, &PD.UploadedAt)
	return PD, err
}

func RemovePerformanceDocument(DocumentID int) error {
	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM performance_documents where document_id = $1`, DocumentID)
	return err
}
