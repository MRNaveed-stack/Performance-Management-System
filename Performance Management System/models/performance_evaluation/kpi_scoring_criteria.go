package performance_evaluation

import (
	"context"
	"performanceManagement/config"
)

type KpiScoringCriteria struct {
	KpiScoringCriteriaID int    `json:"kpi_scoring_criteria_id"`
	KpiID                int    `json:"kpi_id"`
	MinValue             int    `json:"min_value"`
	MaxValue             int    `json:"max_value"`
	Score                int    `json:"score"`
	Description          string `json:"description"`
}

func AddKpiScoringCriteria(KSC KpiScoringCriteria) error {
	_, err := config.DB.Exec(context.Background(),
		`SELECT sp_insert_kpi_scoring_criteria`, KSC.KpiID, KSC.MinValue, KSC.MaxValue, KSC.Score, KSC.Description)
	return err

}

func UpdateKpiScoringCriteria(KSC KpiScoringCriteria) error {

	_, err := config.DB.Exec(context.Background(),
		`SELECT sp_update_kpi_scoring_criteria`, KSC.KpiScoringCriteriaID, KSC.KpiID, KSC.MinValue,
		KSC.MaxValue, KSC.Score, KSC.Description)
	return err
}

func GetKpiCriteria(KpiID int) ([]KpiScoringCriteria, error) {
	var criteria []KpiScoringCriteria
	query := `SELECT * FROM fn_list_kpi_scoring_criteria($1)`
	rows, err := config.DB.Query(context.Background(), query, KpiID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c KpiScoringCriteria
		err := rows.Scan(
			&c.KpiScoringCriteriaID,
			&c.KpiID,
			&c.MinValue,
			&c.MaxValue,
			&c.Score,
			&c.Description,
		)

		if err != nil {
			return nil, err
		}
		criteria = append(criteria, c)
	}
	return criteria, nil
}
