package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type SelfAppraisal struct {
	SelfAppraisalID     int       `json:"self_appraisal_id"`
	PerformanceCycleID  int       `json:"performance_cycle_id"`
	PerformanceCycle    string    `json:"performance_cycle_name"`
	EmployeeID          int       `json:"employee_id"`
	EmployeeName        string    `json:"employee_name"`
	JobDefinition       string    `json:"job_definition"`
	KeyResponsibilities string    `json:"key_responsibilities"`
	Accomplishments     string    `json:"accomplishments"`
	Challenges          string    `json:"challenges"`
	GoalsNextYear       string    `json:"goals_next_year"`
	ManagerSupport      string    `json:"manager_support"`
	DevelopmentNeeds    string    `json:"development_needs"`
	EmployeeProjectID   *int      `json:"employee_project_id,omitempty"`
	ProjectName         *string   `json:"project_name,omitempty"`
	SubmittedAt         time.Time `json:"submitted_at"`
}

func InsertSelfAppraisal(SA SelfAppraisal) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_insert_self_appraisal($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		SA.PerformanceCycleID, SA.EmployeeID, SA.JobDefinition, SA.KeyResponsibilities,
		SA.Accomplishments, SA.Challenges, SA.GoalsNextYear, SA.ManagerSupport,
		SA.DevelopmentNeeds, SA.EmployeeProjectID)
	return err
}

func UpdateSelfAppraisal(SA SelfAppraisal) error {
	_, err := config.DB.Exec(context.Background(),
		`CALL sp_update_self_appraisal($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		SA.SelfAppraisalID, SA.JobDefinition, SA.KeyResponsibilities, SA.Accomplishments, SA.Challenges,
		SA.GoalsNextYear, SA.ManagerSupport, SA.DevelopmentNeeds, SA.EmployeeProjectID)
	return err
}

func GetSelfAppraisal() ([]SelfAppraisal, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT * FROM fn_list_self_appraisals()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var selfappraisal []SelfAppraisal
	for rows.Next() {
		var SA SelfAppraisal
		err := rows.Scan(&SA.SelfAppraisalID, &SA.PerformanceCycleID, &SA.PerformanceCycle,
			&SA.EmployeeID, &SA.EmployeeName,
			&SA.JobDefinition, &SA.KeyResponsibilities, &SA.Accomplishments,
			&SA.Challenges, &SA.Challenges, &SA.GoalsNextYear,
			&SA.ManagerSupport, &SA.DevelopmentNeeds, &SA.EmployeeProjectID,
			&SA.ProjectName, &SA.SubmittedAt)
		if err != nil {
			return nil, err
		}
		selfappraisal = append(selfappraisal, SA)
	}
	return selfappraisal, nil
}

// Get by ID
func GetSelfAppraisalByID(id int) (SelfAppraisal, error) {
	var sa SelfAppraisal
	err := config.DB.QueryRow(context.Background(),
		`SELECT * FROM fn_get_self_appraisal_by_id($1)`, id).
		Scan(
			&sa.SelfAppraisalID, &sa.PerformanceCycleID,
			&sa.PerformanceCycle, &sa.EmployeeID,
			&sa.EmployeeName, &sa.JobDefinition,
			&sa.KeyResponsibilities, &sa.Accomplishments,
			&sa.Challenges, &sa.GoalsNextYear,
			&sa.ManagerSupport, &sa.DevelopmentNeeds,
			&sa.EmployeeProjectID, &sa.ProjectName, &sa.SubmittedAt,
		)
	return sa, err
}
