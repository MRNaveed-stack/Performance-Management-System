package attendance

import (
	"context"
	"performanceManagement/config"
	"time"
)

type OvertimeRecord struct {
	EmployeeID      int     `json:"employee_id"`
	DepartmentID    int     `json:"employee_department_id"`
	OvertimeDate    string  `json:"overtime_date"`
	TimeIn          string  `json:"time_in"`
	TimeOut         string  `json:"time_out"`
	TotalHours      float64 `json:"total_hours"`
	OvertimeHours   float64 `json:"overtime_hours"`
	HourlyRate      float64 `json:"hourly_rate"`
	OvertimePayment float64 `json:"overtime_payment"`
}

// SaveOvertimeRecord inserts or updates an overtime record via stored procedure
func SaveOvertimeRecord(
	employeeID int,
	departmentID int,
	date time.Time,
	timeIn time.Time,
	timeOut time.Time,
	totalHours float64,
	overtimeHours float64,
	hourlyRate float64,
	overtimePayment float64,
) error {

	// time_in and time_out are always saved as HH:MM:SS
	_, err := config.DB.Exec(
		context.Background(),
		`
		CALL sp_save_overtime_record(
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
		`,
		employeeID,
		departmentID,
		date, // DATE column in DB
		timeIn.Format("15:04:05"),
		timeOut.Format("15:04:05"),
		totalHours,
		overtimeHours,
		hourlyRate,
		overtimePayment,
	)
	return err
}
