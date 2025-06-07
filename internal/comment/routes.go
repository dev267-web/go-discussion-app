// routes.go 
package comment

import (
    "database/sql"

    "github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB) {
    repo := NewRepository(db)
    svc := NewService(repo)
    ctr := NewController(svc)

    rg.POST("/discussions/:id/comments", ctr.Create)
    rg.GET("/discussions/:id/comments", ctr.List)
}
