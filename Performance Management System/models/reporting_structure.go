package models

import (
	"context"
	"fmt"
	"performanceManagement/config"
)

type ReportingStructure struct {
	EmployeeID         int    `json:"employee_id"`
	ReportingManagerID *int   `json:"reporting_manager_id,omitempty"`
	ReportingManager   string `json:"reporting_manager_name,omitempty"`
	TeamID             *int   `json:"employee_team_id,omitempty"`
	TeamName           string `json:"team_name,omitempty"`
}

type Team struct {
	TeamID   int    `json:"team_id"`
	TeamName string `json:"team_name"`
}

func FetchAllTeams() ([]Team, error) {
	rows, err := config.DB.Query(context.Background(), "SELECT team_id, team_name FROM teams ORDER BY team_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.TeamID, &t.TeamName); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func GetReportingStructureByID(empID int) (*ReportingStructure, error) {
	query := `
		SELECT e.employee_id,
			   e.reporting_manager_id,
			   COALESCE(m.employee_name, '') AS reporting_manager_name,
			   e.employee_team_id,
			   COALESCE(t.team_name, '') AS team_name
		FROM employees e
		LEFT JOIN employees m ON e.reporting_manager_id = m.employee_id
		LEFT JOIN teams t ON e.employee_team_id = t.team_id
		WHERE e.employee_id = $1
	`

	var rs ReportingStructure
	err := config.DB.QueryRow(context.Background(), query, empID).Scan(
		&rs.EmployeeID,
		&rs.ReportingManagerID,
		&rs.ReportingManager,
		&rs.TeamID,
		&rs.TeamName,
	)

	if err != nil {
		return nil, err
	}
	return &rs, nil
}

func UpdateReportingStructure(empID, managerID, teamID int) error {
	query := `
		UPDATE employees
		SET reporting_manager_id = $1,
			employee_team_id = $2
		WHERE employee_id = $3
	`
	_, err := config.DB.Exec(context.Background(), query, managerID, teamID, empID)
	return err
}

func SetReportingStructureByNames(empID int, managerName, teamName string) error {
	var managerID, teamID int

	// Lookup manager ID
	if err := config.DB.QueryRow(context.Background(),
		`SELECT employee_id FROM employees WHERE employee_name = $1`, managerName).Scan(&managerID); err != nil {
		return fmt.Errorf("manager not found: %v", err)
	}

	// Lookup team ID
	if err := config.DB.QueryRow(context.Background(),
		`SELECT team_id FROM teams WHERE team_name = $1`, teamName).Scan(&teamID); err != nil {
		return fmt.Errorf("team not found: %v", err)
	}

	// Update structure
	return UpdateReportingStructure(empID, managerID, teamID)
}

type ReportingManager struct {
	EmployeeID int    `json:"employee_id"`
	FullName   string `json:"employee_name"`
}

func FetchReportingManagers() ([]ReportingManager, error) {
	query := `
		SELECT e.employee_id, e.employee_name
		FROM employees e
		WHERE e.employee_designation IN ('HR', 'GM')
		ORDER BY e.employee_name;
	`

	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managers []ReportingManager
	for rows.Next() {
		var rm ReportingManager
		if err := rows.Scan(&rm.EmployeeID, &rm.FullName); err != nil {
			return nil, err
		}
		managers = append(managers, rm)
	}

	return managers, nil
}
