package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type Kpi struct {
	KpiID                    int       `json:"kpi_id"`
	KpiCategoryID            int       `json:"kpi_category_id"`
	PerformanceCycleStatusID int       `json:"performance_cycle_status_id"`
	KpiName                  string    `json:"kpi_name"`
	Description              string    `json:"description"`
	Weightage                int       `json:"weightage"`
	CreatedBy                int       `json:"created_by"`
	CreatedAt                time.Time `json:"created_at"`
}

func AddKpi(kpi Kpi) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_insert_kpi($1,$2,$3,$4,$5,$6,$7)`,
		kpi.KpiCategoryID, kpi.PerformanceCycleStatusID, kpi.KpiName, kpi.Description, kpi.Weightage,
		kpi.CreatedBy, kpi.CreatedAt)
	return err
}

func UpdateKpi(kpi Kpi) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_update_kpi($1,$2,$3,$4,$5,$6,$7)`,
		kpi.KpiID, kpi.PerformanceCycleStatusID, kpi.KpiName,
		kpi.Description, kpi.Weightage, kpi.CreatedBy)
	return err
}

func ListKpi(KpiCategoryID *int) ([]Kpi, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * from fn_list_kpis($1)`, KpiCategoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kpi []Kpi
	for rows.Next() {

		var k Kpi
		err := rows.Scan(&k.KpiID, &k.KpiCategoryID, &k.PerformanceCycleStatusID, &k.KpiName, &k.Description,
			&k.Weightage, &k.CreatedBy, &k.CreatedAt)
		if err != nil {
			return nil, err
		}
		kpi = append(kpi, k)
	}

	return kpi, nil
}
