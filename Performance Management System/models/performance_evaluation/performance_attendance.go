package performance_evaluation

import (
	"context"
	"performanceManagement/config"
	"time"
)

type PerformanceAttendance struct {
	AttendanceID       int       `json:"attendance_id"`
	EmployeeID         int       `json:"employee_id"`
	PerformanceCycleID int       `json:"performance_cycle_id"`
	OnTimeLogins       float64   `json:"on_time_logins"`
	HoursCompletion    float64   `json:"hours_completion"`
	CalculatedScore    float64   `json:"calculated_score"`
	ReviewedBy         int       `json:"reviewed_by"`
	ReviewedAt         time.Time `json:"reviewed_at"`
}

func InsertPerformanceAttendance(a PerformanceAttendance) (int, error) {
	var id int
	err := config.DB.QueryRow(context.Background(),
		"SELECT sp_insert_performance_attendance($1,$2,$3,$4,$5)",
		a.EmployeeID, a.PerformanceCycleID, a.OnTimeLogins, a.HoursCompletion, a.ReviewedBy,
	).Scan(&id)
	return id, err
}

func UpdatePerformanceAttendance(a PerformanceAttendance) error {
	_, err := config.DB.Exec(context.Background(),
		"CALL sp_update_performance_attendance($1,$2,$3,$4)",
		a.AttendanceID, a.OnTimeLogins, a.HoursCompletion, a.ReviewedBy,
	)
	return err
}

func ListPerformanceAttendance(employeeID int) ([]PerformanceAttendance, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT * FROM fn_list_performance_attendance($1)", employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PerformanceAttendance
	for rows.Next() {
		var pa PerformanceAttendance
		if err := rows.Scan(&pa.AttendanceID, &pa.PerformanceCycleID, &pa.OnTimeLogins, &pa.HoursCompletion,
			&pa.CalculatedScore, &pa.ReviewedBy, &pa.ReviewedAt); err != nil {
			return nil, err
		}
		pa.EmployeeID = employeeID
		list = append(list, pa)
	}
	return list, nil
}
