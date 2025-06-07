// routes.go 
package tag

import (
    "database/sql"

    "github.com/gin-gonic/gin"
)

// RegisterRoutes mounts the /tags endpoint onto the given router group.
// This should be called on your protected router group in main.go.
func RegisterRoutes(rg *gin.RouterGroup, dbConn *sql.DB) {
    repo := NewRepository(dbConn)
    svc := NewService(repo)
    ctr := NewController(svc)

    rg.GET("/tags", ctr.ListHandler)
}
