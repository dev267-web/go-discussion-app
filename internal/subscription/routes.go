// routes.go 
package subscription

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewSubscriptionController(service)

	rg.POST("/discussions/:id/subscribe", controller.Subscribe)
	rg.DELETE("/discussions/:id/unsubscribe", controller.Unsubscribe)
	rg.POST("/discussions/:id/notify", controller.Notify)
}
