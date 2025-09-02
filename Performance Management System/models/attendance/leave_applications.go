package attendance

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)

type LeaveApplication struct {
	LeaveApplicationID int       `json:"leave_application_id"`
	EmployeeID         int       `json:"employee_id"`
	LeaveTypeID        int       `json:"leave_type_id"`
	DepartmentID       int       `json:"department_id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	FromDate           time.Time `json:"from_date"`
	ToDate             time.Time `json:"to_date"`
	NumberOfLeaves     int       `json:"number_of_leaves"`
	AttachedDocument   string    `json:"attached_document"`
	CreatedDate        time.Time `json:"created_date"`
	CurrentStatusID    int       `json:"current_status_id"`
}

type NullInt32 struct {
	sql.NullInt32
}

func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.Itoa(int(ni.Int32))), nil
}

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

type LeaveApplicationRequest struct {
	LeaveApplicationID      int            `json:"leave_application_id"`
	EmployeeID              int            `json:"employee_id"`
	LeaveTypeID             int            `json:"leave_type_id"`
	DepartmentID            NullInt32      `json:"department_id"`
	Title                   string         `json:"title"`
	Description             string         `json:"description"`
	FromDate                time.Time      `json:"from_date"`
	ToDate                  time.Time      `json:"to_date"`
	AttachedDocument        sql.NullString `json:"attached_document"`
	CurrentStatusID         sql.NullInt32  `json:"current_status_id"`
	ReportingManagerStatus  string         `json:"reporting_manager_status"`
	ReportingManagerRemarks sql.NullString `json:"reporting_manager_remarks"`
	HRStatus                string         `json:"hr_status"`
	HRRemarks               sql.NullString `json:"hr_remarks"`
}

type LeaveApplicationRequest1 struct {
	LeaveApplicationID      int            `json:"leave_application_id"`
	EmployeeID              int            `json:"employee_id"`
	LeaveTypeID             int            `json:"leave_type_id"`
	DepartmentID            *int           `json:"department_id"`
	Title                   string         `json:"title"`
	Description             string         `json:"description"`
	FromDate                time.Time      `json:"from_date"`
	ToDate                  time.Time      `json:"to_date"`
	NumberOfLeaves          sql.NullInt32  `json:"number_of_leaves"`
	AttachedDocument        sql.NullString `json:"attached_document"`
	CreatedDate             time.Time      `json:"created_date"`
	CurrentStatusID         sql.NullInt32  `json:"current_status_id"`
	ReportingManagerStatus  string         `json:"reporting_manager_status"`
	ReportingManagerRemarks sql.NullString `json:"reporting_manager_remarks"`
	HRStatus                string         `json:"hr_status"`
	HRRemarks               sql.NullString `json:"hr_remarks"`
}
