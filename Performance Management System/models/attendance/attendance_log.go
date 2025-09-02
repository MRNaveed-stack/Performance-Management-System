package attendance

import "time"

type PunchLog struct {
	PunchID       int       `json:"punch_id"`
	EmployeeID    int       `json:"employee_id"`
	PunchDateTime time.Time `json:"punch_date_time"`
	PunchTypeID   int       `json:"punch_type_id"`
}

type AttendanceRecord struct {
	AttendanceID int        `json:"attendance_id"`
	EmployeeID   int        `json:"employee_id"`
	DepartmentID int        `json:"department_id"`
	Date         time.Time  `json:"date"`
	TimeIn       *time.Time `json:"time_in"`
	TimeOut      *time.Time `json:"time_out"`
	TotalHours   float64    `json:"total_hours"`
	Status       string     `json:"status"`
}

type AttendanceRecordResponse struct {
	AttendanceID   int       `json:"attendance_id"`
	EmployeeID     int       `json:"employee_id"`
	EmployeeName   string    `json:"employee_name"`
	DepartmentName string    `json:"department_name"`
	Date           time.Time `json:"date"`
	TimeIn         string    `json:"time_in,omitempty"`
	TimeOut        string    `json:"time_out,omitempty"`
	TotalHours     float64   `json:"total_hours"`
	Status         string    `json:"status"`
}
