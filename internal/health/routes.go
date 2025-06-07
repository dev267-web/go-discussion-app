// routes.go 
package health

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	service := NewHealthService(db)
	controller := NewHealthController(service)

	r.GET("/health", controller.HandleHealthCheck)
}
