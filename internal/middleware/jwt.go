// jwt.go 
package middleware

import "go-discussion-app/internal/auth"
import "github.com/gin-gonic/gin"

// JWTAuth is the shared alias for auth.JWTAuthMiddleware
func JWTAuth() gin.HandlerFunc {
  return auth.JWTAuthMiddleware()
}
