package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type PerformanceWorkflow struct {
	PerformanceWorkflowID    int       `json:"performance_workflow_id"`
	PerformanceCycleID       int       `json:"performance_cycle_id"`
	EmployeeID               int       `json:"employee_id"`
	CurrentStage             string    `json:"current_stage"`
	PerformanceCycleStatusID int       `json:"performance_cycle_status_id"`
	StartedAt                time.Time `json:"started_at"`
	CompletedAt              time.Time `json:"completed_at"`
}

func InsertPeroformanceWorkflow(PW PerformanceWorkflow) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_insert_performance_workflow($1,$2,$3,$4)`,
		PW.PerformanceCycleID, PW.EmployeeID, PW.CurrentStage, PW.PerformanceCycleStatusID)
	return err
}

func UpdatePerformanceWorkflow(PW PerformanceWorkflow) error {
	_, err := config.DB.Exec(context.Background(),
		`SELECT sp_update_performance_workflow($1,$2,$3,$4)`,
		PW.PerformanceWorkflowID, PW.CurrentStage, PW.PerformanceCycleStatusID, PW.CompletedAt)
	return err
}

func GetPerformanceWorkflow(EmployeeID int) ([]PerformanceWorkflow, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * from fn_list_performance_workflows_by_employee($1)`, EmployeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performanceworkflow []PerformanceWorkflow
	for rows.Next() {
		var PW PerformanceWorkflow
		err := rows.Scan(&PW.PerformanceWorkflowID, &PW.PerformanceCycleID, PW.EmployeeID,
			&PW.CurrentStage, &PW.PerformanceCycleStatusID, &PW.StartedAt, &PW.CompletedAt)
		if err != nil {
			return nil, err
		}
		performanceworkflow = append(performanceworkflow, PW)
	}
	return performanceworkflow, nil
}
