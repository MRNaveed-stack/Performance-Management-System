package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type KpiTarget struct {
	KpiTargetID     int       `json:"kpi_target_id"`
	KpiAssignmentID int       `json:"kpi_assignment_id"`
	TargetValue     string    `json:"target_value"`
	BenchMarkValue  string    `json:"benchmark_value"`
	Criteria        string    `json:"criteria"`
	ReviewCycle     string    `json:"review_cycle"`
	CreatedAt       time.Time `json:"created_at"`
}

func InsertKpiTarget(KT KpiTarget) error {
	_, err := config.DB.Exec(context.Background(),
		`SELECT sp_insert_kpi_target($1,$2,$3,$4,$5)`,
		KT.KpiAssignmentID, KT.TargetValue, KT.BenchMarkValue, KT.Criteria,
		KT.ReviewCycle)
	return err
}

func UpdateKpiTarget(KT KpiTarget) error {
	_, err := config.DB.Exec(context.Background(),
		`SELECT sp_update_kpi_target($1,$2,$3,$4,$5)`,
		KT.KpiTargetID, KT.TargetValue, KT.BenchMarkValue, KT.Criteria, KT.ReviewCycle)
	return err
}

func GetKpiTarget() ([]KpiTarget, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * FROM fn_list_kpi_targets()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kpitarget []KpiTarget
	for rows.Next() {
		var t KpiTarget
		err := rows.Scan(&t.KpiTargetID, &t.KpiAssignmentID, &t.TargetValue,
			&t.BenchMarkValue, &t.Criteria, &t.ReviewCycle, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		kpitarget = append(kpitarget, t)
	}
	return kpitarget, nil
}
