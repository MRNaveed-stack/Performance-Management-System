package attendance

import (
	"context"
	"database/sql"
	"net/http"
	"performanceManagement/config"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAbsenceTracking(c *gin.Context) {
	date := c.Query("date")
	dept := c.Query("department")
	search := c.Query("search")

	// default to today
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		// accept common formats (YYYY-MM-DD or M/D/YYYY or MM/DD/YYYY)
		var parsed time.Time
		var err error
		layouts := []string{"2006-01-02", "1/2/2006", "01/02/2006"}
		for _, l := range layouts {
			parsed, err = time.Parse(l, date)
			if err == nil {
				date = parsed.Format("2006-01-02")
				break
			}
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format; use YYYY-MM-DD or M/D/YYYY"})
			return
		}
	}

	if dept == "" {
		dept = ""
	}
	if search == "" {
		search = ""
	}

	// include people who have no attendance record OR have an attendance record marked 'Absent'
	query := `
SELECT 
    e.employee_id,
    e.employee_name,
    d.department_name,
    COALESCE(lt.leave_type_name, 'No Leave Request') AS leave_type,
    -- two-step approval logic for leave status
    COALESCE(
        CASE 
            WHEN la.reporting_manager_status != 'Approved' 
                THEN la.reporting_manager_status
            WHEN la.reporting_manager_status = 'Approved' 
                THEN la.hr_status
            ELSE 'Absent Without Leave'
        END,
        'Absent Without Leave'
    ) AS status
FROM employees e
JOIN departments d 
    ON e.employee_department_id = d.department_id
LEFT JOIN attendance_records ar
    ON ar.employee_id = e.employee_id
    AND ar.date = $1::date
LEFT JOIN leave_applications la
    ON la.employee_id = e.employee_id
    AND $1::date BETWEEN la.from_date AND la.to_date
LEFT JOIN leave_type lt 
    ON la.leave_type_id = lt.leave_type_id
WHERE 
    (ar.status IS NULL OR ar.status = 'Absent')
    AND ($2 = '' OR d.department_name = $2)
    AND (
        $3 = '' 
        OR e.employee_name ILIKE '%' || $3 || '%' 
        OR e.employee_id::TEXT ILIKE '%' || $3 || '%'
    )
ORDER BY e.employee_id;
`

	rows, err := config.DB.Query(context.Background(), query, date, dept, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		var empID int
		var empName, deptName, leaveType, status sql.NullString

		if err := rows.Scan(&empID, &empName, &deptName, &leaveType, &status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		records = append(records, map[string]interface{}{
			"employee_id":     empID,
			"employee_name":   empName.String,
			"department_name": deptName.String,
			"leave_type":      leaveType.String,
			"status":          status.String,
		})
	}

	c.JSON(http.StatusOK, records)
}
