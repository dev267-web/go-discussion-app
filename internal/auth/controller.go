// controller.go 
package auth

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "go-discussion-app/pkg/logger"
)

type AuthController struct {
    svc *AuthService
}

func NewController(svc *AuthService) *AuthController {
    return &AuthController{svc: svc}
}

func (ctr *AuthController) RegisterHandler(c *gin.Context) {
    var dto RegisterDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        logger.Warnf("register binding error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    id, err := ctr.svc.Register(c.Request.Context(), &dto)
    if err != nil {
        if err == ErrUserExists {
            c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
        } else {
            logger.Errorf("register error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
        }
        return
    }
    c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ctr *AuthController) LoginHandler(c *gin.Context) {
    var dto LoginDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        logger.Warnf("login binding error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    token, err := ctr.svc.Login(c.Request.Context(), &dto)
    if err != nil {
        if err == ErrInvalidCredentials {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong email or password"})
        } else {
            logger.Errorf("login error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}
