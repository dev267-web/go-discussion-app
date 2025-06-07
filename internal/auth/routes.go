// routes.go 
package auth

import (
    "database/sql"

    "github.com/gin-gonic/gin"
    "go-discussion-app/internal/user"
)

// RegisterRoutes mounts /auth/register and /auth/login.
// Pass router, DB connection, and the JWT secret (if you want to use it in middleware).
func RegisterRoutes(router *gin.Engine, dbConn *sql.DB) {
    userRepo := user.NewRepository(dbConn)
    svc := NewService(userRepo)
    ctr := NewController(svc)

    grp := router.Group("/auth")
    grp.POST("/register", ctr.RegisterHandler)
    grp.POST("/login", ctr.LoginHandler)
}
