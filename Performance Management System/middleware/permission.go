package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"performanceManagement/config"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")
		role := c.GetString("role")

		// HR Manager full access
		if strings.EqualFold(role, "HR Manager") || strings.EqualFold(role, "hr") {
			c.Next()
			return
		}

		endpoint := c.FullPath()
		method := c.Request.Method
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var count int

		query := `
SELECT COALESCE(SUM(c), 0) AS total_count
FROM (
    SELECT COUNT(*) AS c
    FROM api_permissions ap
    JOIN user_roles ur ON ur.role_id = ap.role_id
    WHERE ur.user_id = $1 AND ap.method = $2 AND ap.endpoint = $3

    UNION ALL

    SELECT COUNT(*) AS c
    FROM user_permissions up
    JOIN permissions p ON p.id = up.permission_id
    WHERE up.user_id = $1 AND p.name = $2 || ' ' || $3
) AS combined
`

		err := config.DB.QueryRow(ctx, query, userID, method, endpoint).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
	}
}
