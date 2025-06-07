// controller.go 
package user

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go-discussion-app/pkg/logger"
    //"go-discussion-app/models"
)

type UserController struct {
    svc *UserService
}

func NewController(svc *UserService) *UserController {
    return &UserController{svc: svc}
}

// GetProfile handles GET /users/:id
func (ctr *UserController) GetProfile(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    user, err := ctr.svc.GetByID(c.Request.Context(), id)
    if err != nil {
        if err == ErrUserNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        } else {
            logger.Errorf("GetProfile error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
        }
        return
    }

    // Hide password hash in response
    user.PasswordHash = ""
    c.JSON(http.StatusOK, user)
}

// UpdateProfile handles PUT /users/:id
func (ctr *UserController) UpdateProfile(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    var dto UpdateUserDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    if err := dto.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updated, err := ctr.svc.Update(c.Request.Context(), id, &dto)
    if err != nil {
        switch err {
        case ErrUserNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        default:
            logger.Errorf("UpdateProfile error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
        }
        return
    }

    updated.PasswordHash = ""
    c.JSON(http.StatusOK, updated)
}

// DeleteProfile handles DELETE /users/:id
func (ctr *UserController) DeleteProfile(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    if err := ctr.svc.Delete(c.Request.Context(), id); err != nil {
        switch err {
        case ErrUserNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        default:
            logger.Errorf("DeleteProfile error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
        }
        return
    }
    c.Status(http.StatusNoContent)
}
