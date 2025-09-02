package performance_evaluation

import (
	"context"
	"performanceManagement/config"
)

type EmployeeKpiAssignment struct {
	EmployeeKpiAssignmentID int    `json:"employee_kpi_assignment_id"`
	PerformanceCycleID      int    `json:"performance_cycle_id"`
	EmployeeID              int    `json:"employee_id"`
	ReportingManagerID      int    `json:"reporting_manager_id"`
	KpiID                   int    `json:"kpi_id"`
	TargeValue              string `json:"target_value"`
}

func InsertEmployeeKpiAssignment(EKA EmployeeKpiAssignment) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_insert_employee_kpi_assignment($1,$2,$3,$4,$5)`,
		EKA.PerformanceCycleID, EKA.EmployeeID, EKA.ReportingManagerID, EKA.KpiID, EKA.TargeValue)
	return err
}

func UpdateEmployeeKpiAssignment(EKA EmployeeKpiAssignment) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_update_employee_kpi_assignment($1,$2,$3,$4,$5)`,
		EKA.PerformanceCycleID, EKA.EmployeeID, EKA.ReportingManagerID, EKA.KpiID, EKA.TargeValue)
	return err
}

func GetEmployeeKpiAssignment(employeeID int) ([]EmployeeKpiAssignment, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * FROM fn_list_employee_kpi_assignments($1)`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []EmployeeKpiAssignment
	for rows.Next() {
		var R EmployeeKpiAssignment
		err := rows.Scan(&R.EmployeeKpiAssignmentID, &R.PerformanceCycleID, &R.EmployeeID,
			&R.ReportingManagerID, &R.KpiID, &R.TargeValue)
		if err != nil {
			return nil, err
		}
		result = append(result, R)
	}
	return result, nil
}
