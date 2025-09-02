package models

import (
	"context"
	"fmt"
	"performanceManagement/config"
	"time"
)

type Certification struct {
	CertificateID       int       `json:"certificate_id"`
	EmployeeID          int       `json:"employee_id"`
	CertificationName   string    `json:"certification_name"`
	IssuingOrganization string    `json:"issuing_organization"`
	DateAchieved        time.Time `json:"date_achieved"`
}

func CreateCertification(cert Certification) error {
	query := `
		INSERT INTO certifications (
			employee_id, certification_name, issuing_organization, date_achieved
		) VALUES ($1, $2, $3, $4)
	`
	_, err := config.DB.Exec(context.Background(), query,
		cert.EmployeeID,
		cert.CertificationName,
		cert.IssuingOrganization,
		cert.DateAchieved,
	)
	if err != nil {
		return fmt.Errorf("DB insert failed: %w", err)
	}
	return nil
}
func GetCertificationsByEmployee(empID int) ([]Certification, error) {
	query := `
		SELECT certificate_id, employee_id, certification_name,
			   issuing_organization, date_achieved
		FROM certifications
		WHERE employee_id = $1
	`

	rows, err := config.DB.Query(context.Background(), query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certs []Certification
	for rows.Next() {
		var c Certification
		err := rows.Scan(
			&c.CertificateID,
			&c.EmployeeID,
			&c.CertificationName,
			&c.IssuingOrganization,
			&c.DateAchieved,
		)
		if err != nil {
			return nil, err
		}
		certs = append(certs, c)
	}
	return certs, nil
}

func DeleteCertification(certID int) error {
	_, err := config.DB.Exec(context.Background(), `DELETE FROM certifications WHERE certificate_id = $1`, certID)
	return err
}
