package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type PerformanceReview struct {
	PerformanceReviewID     int       `json:"performance_review_id"`
	EmployeeKpiAssignmentID int       `json:"employee_kpi_assignment_id"`
	PerformanceReviewTypeID int       `json:"performance_review_type_id"`
	GivenScore              int       `json:"given_score"`
	Comments                string    `jon:"comments"`
	ReviewedBy              int       `json:"reviewed_by"`
	ReviewerName            string    `json:"reviewer_name"`
	ReviewedAt              time.Time `json:"reviewed_at"`
}

func InsertPerformanceReview(PR PerformanceReview) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_insert_performance_review($1,$2,$3,$4,$5)`,
		PR.EmployeeKpiAssignmentID, PR.PerformanceReviewTypeID, PR.GivenScore, PR.Comments, PR.ReviewedBy)
	return err
}

func UpdatePerformanceReview(PR PerformanceReview) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_update_performance_review($1,$2,$3,$4,$5,$6)`,
		PR.EmployeeKpiAssignmentID, PR.PerformanceReviewTypeID, PR.GivenScore,
		PR.Comments, PR.ReviewedBy, PR.ReviewedAt)
	return err
}

func GetPerformanceReviewByID(reviewID int) (PerformanceReview, error) {
	row := config.DB.QueryRow(context.Background(),
		"SELECT * FROM fn_get_performance_review_by_id($1)", reviewID)

	var pr PerformanceReview
	err := row.Scan(
		&pr.PerformanceReviewID,
		&pr.EmployeeKpiAssignmentID,
		&pr.PerformanceReviewTypeID,
		&pr.GivenScore,
		&pr.Comments,
		&pr.ReviewedBy,
		&pr.ReviewerName,
		&pr.ReviewedAt,
	)
	return pr, err
}

func GetPerformanceReview() ([]PerformanceReview, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT * FROM fn_list_performance_reviews()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performancereview []PerformanceReview
	for rows.Next() {

		var pr PerformanceReview
		err := rows.Scan(
			&pr.PerformanceReviewID,
			&pr.EmployeeKpiAssignmentID,
			&pr.PerformanceReviewTypeID,
			&pr.GivenScore,
			&pr.Comments,
			&pr.ReviewedBy,
			&pr.ReviewerName,
			&pr.ReviewedAt,
		)
		if err != nil {
			return nil, err
		}
		performancereview = append(performancereview, pr)
	}

	return performancereview, nil
}
