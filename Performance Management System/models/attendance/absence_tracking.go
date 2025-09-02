package attendance

type AbsenceRecord struct {
	EmployeeID string `json:"employee_id"`
	FullName   string `json:"employee_name"`
	Department string `json:"departments"`
	LeaveType  string `json:"leave_type"`
	Status     string `json:"status"`
}
