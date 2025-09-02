package attendance

import (
	"net/http"

	"performanceManagement/config"

	"github.com/gin-gonic/gin"
)

type PunchRecord struct {
	DeviceSerial string  `json:"device_serial_num"`
	EnrollID     int     `json:"enroll_id"`
	Event        int     `json:"event"`
	InOut        int     `json:"intOut"`
	Mode         int     `json:"mode"`
	RecordsTime  string  `json:"records_time"`
	Temperature  float64 `json:"temperature"`
}

type Payload struct {
	Records []PunchRecord `json:"records"`
}

func PunchHandler(c *gin.Context) {
	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, rec := range payload.Records {
		punchType := 1 // default OUT
		if rec.InOut == 0 {
			punchType = 0 // IN
		}

		_, err := config.DB.Exec(
			c,
			`INSERT INTO punch_logs (employee_id, punch_date_time, punch_type_id)
			 VALUES ($1, $2, $3)
			 ON CONFLICT DO NOTHING`,
			rec.EnrollID, rec.RecordsTime, punchType+1, // punch_type_id 1=IN, 2=OUT
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Punch logs inserted successfully"})
}
