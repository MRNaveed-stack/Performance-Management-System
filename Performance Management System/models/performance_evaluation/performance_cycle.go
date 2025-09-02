package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type PerformanceCycle struct {
	PerformanceCycleID          int       `json:"performance_cycle_id"`
	PerformanceCycleName        string    `json:"performance_cycle_name"`
	DepartmentID                int       `json:"department_id"`
	Description                 string    `json:"description"`
	PerformanceCycleFrequencyId int       `json:"performance_cycle_frequency_id"`
	StartDate                   time.Time `json:"start_date"`
	EndDate                     time.Time `json:"end_date"`
	PerformanceCycleStatusID    int       `json:"performance_cycle_status_id"`
	CreatedBy                   int       `json:"created_by"`
	CreatedAt                   time.Time `json:"created_at"`
}

func InsertPerformanceCycle(cycle PerformanceCycle) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_insert_performance_cycle($1,$2,$3,$4,$5,$6,$7,$8)`,
		cycle.PerformanceCycleName, cycle.DepartmentID, cycle.Description, cycle.PerformanceCycleFrequencyId,
		cycle.StartDate, cycle.EndDate, cycle.PerformanceCycleStatusID, cycle.CreatedBy)
	return err
}

func UpdatePerformanceCycle(cycle PerformanceCycle) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_update_performance_cycle($1,$2,$3,$4,$5,$6,$7)`,
		cycle.PerformanceCycleName, cycle.DepartmentID, cycle.Description,
		cycle.PerformanceCycleFrequencyId, cycle.StartDate, cycle.EndDate, cycle.PerformanceCycleStatusID)
	return err
}

func GetPerformanceCycle() ([]PerformanceCycle, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * FROM fn_list_performance_cycles()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cycles []PerformanceCycle
	for rows.Next() {
		var C PerformanceCycle
		err := rows.Scan(&C.PerformanceCycleID, &C.PerformanceCycleName, &C.DepartmentID, &C.Description,
			&C.PerformanceCycleFrequencyId, &C.StartDate, &C.EndDate, &C.PerformanceCycleStatusID,
			&C.CreatedBy, &C.CreatedAt)
		if err != nil {
			return nil, err
		}
		cycles = append(cycles, C)
	}
	return cycles, nil
}
