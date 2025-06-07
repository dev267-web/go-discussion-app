// middleware.go 
package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "go-discussion-app/pkg/jwtutil"
)

// JWTAuthMiddleware enforces “Bearer <token>” and sets “userID” in context.
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
            c.Abort()
            return
        }
        uid, err := jwtutil.ExtractUserID(parts[1])
        if err != nil {
            if err == jwtutil.ErrTokenExpired {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
            } else {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            }
            c.Abort()
            return
        }
        c.Set("userID", uid)
        c.Next()
    }
}

// GetUserID retrieves the authenticated user’s ID from context.
func GetUserID(c *gin.Context) (int, bool) {
    raw, exists := c.Get("userID")
    if !exists {
        return 0, false
    }
    uid, ok := raw.(int)
    return uid, ok
}
