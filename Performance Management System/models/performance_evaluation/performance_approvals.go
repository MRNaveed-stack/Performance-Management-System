package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type PerformanceApproval struct {
	ApprovalID     int       `json:"approval_id"`
	WorkFlowID     int       `json:"workflow_id"`
	ApproverID     int       `json:"approver_id"`
	ApproverRole   string    `json:"approver_role"`
	ApprovalStatus string    `json:"approval_status"`
	Comments       string    `json:"comments"`
	ApprovedAt     time.Time `json:"approved_at"`
}

func InsertPerformanceApproval(PA PerformanceApproval) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_insert_performance_approval($1,$2,$3,$4,$5)`,
		PA.WorkFlowID, PA.ApproverID, PA.ApproverRole, PA.ApprovalStatus, PA.Comments)
	return err
}

func UpdatePerformanceApproval(PA PerformanceApproval) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_update_performance_approval($1,$2,$3)`,
		PA.ApprovalID, PA.ApprovalStatus, PA.Comments)
	return err
}

func GetPerformanceApproval(WorkFlowID int) ([]PerformanceApproval, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * FROM fn_list_performance_approvals($1)`, WorkFlowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performanceapproval []PerformanceApproval

	for rows.Next() {
		var PA PerformanceApproval
		err := rows.Scan(&PA.ApprovalID, &PA.WorkFlowID, &PA.ApproverID, &PA.ApprovalStatus, &PA.Comments, &PA.ApprovedAt)

		if err != nil {
			return nil, err
		}
		performanceapproval = append(performanceapproval, PA)
	}
	return performanceapproval, nil
}
