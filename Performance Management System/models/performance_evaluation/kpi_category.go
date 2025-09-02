package performance_evaluation

import (
	"context"
	"performanceManagement/config"
)

type KpiCategory struct {
	CategoryID          int    `json:"kpi_category_id"`
	CategoryName        string `json:"kpi_category_name"`
	CategoryDescription string `json:"kpi_category_description"`
}

func InsertKpiCategory(category KpiCategory) error {
	_, err := config.DB.Exec(context.Background(),

		`SELECT sp_insert_kpi_category($1,$2)`,
		category.CategoryName,
		category.CategoryDescription,
	)
	return err
}

func UpdateKpiCategory(category KpiCategory) error {
	_, err := config.DB.Exec(context.Background(),
		`SELECT  sp_update_kpi_category($1,$2,$3)`,
		category.CategoryID, category.CategoryName, category.CategoryDescription,
	)

	return err
}

func ListKpiCategories() ([]KpiCategory, error) {

	rows, err := config.DB.Query(context.Background(), "SELECT * from fn_list_kpi_categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []KpiCategory
	for rows.Next() {
		var c KpiCategory
		err := rows.Scan(&c.CategoryID, &c.CategoryName, &c.CategoryDescription)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
