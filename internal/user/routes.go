// routes.go 
package user

import (
    "database/sql"

    "github.com/gin-gonic/gin"
)

// RegisterRoutes mounts user/profile endpoints under the protected group.
func RegisterRoutes(rg *gin.RouterGroup, dbConn *sql.DB) {
    repo := NewRepository(dbConn)
    svc := NewService(repo)
    ctr := NewController(svc)

    // All these routes require JWT middleware applied by main.go
    rg.GET("/users/:id", ctr.GetProfile)
    rg.PUT("/users/:id", ctr.UpdateProfile)
    rg.DELETE("/users/:id", ctr.DeleteProfile)
}
