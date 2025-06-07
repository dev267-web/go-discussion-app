// routes.go 
package discussion

import (
    "database/sql"

    "github.com/gin-gonic/gin"
		"go-discussion-app/internal/tag"
)

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB) {
    discRepo := NewRepository(db)

		tagRepo := tag.NewRepository(db)                      // <— new
    svc := NewService(discRepo, tagRepo)
    
    ctr := NewController(svc)

    // standard CRUD
    rg.POST("/discussions", ctr.Create)
    rg.GET("/discussions", ctr.List)
    rg.GET("/discussions/:id", ctr.Get)
    rg.PUT("/discussions/:id", ctr.Update)
    rg.DELETE("/discussions/:id", ctr.Delete)

    // filters & tagging
    rg.GET("/discussions/user/:userId", ctr.ListByUser)
    rg.GET("/discussions/tag/:tag", ctr.ListByTag)
    rg.POST("/discussions/:id/tags", ctr.AddTags)

    // scheduled
    rg.POST("/discussions/schedule", ctr.Schedule)
}
