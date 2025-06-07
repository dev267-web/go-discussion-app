// controller.go 
package comment

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go-discussion-app/pkg/logger"
    "go-discussion-app/internal/auth"
)

type Controller struct {
    svc Service
}

func NewController(svc Service) *Controller {
    return &Controller{svc: svc}
}

// POST /discussions/:id/comments
func (ctr *Controller) Create(c *gin.Context) {
    // Parse discussion ID
    discID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discussion ID"})
        return
    }

    // Bind and validate DTO
    var dto CreateCommentDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    if err := dto.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get userID from context (middleware.JWTAuth must have set it)
    userID, ok := auth.GetUserID(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
        return
    }

    // Call service
    commentID, err := ctr.svc.AddComment(c.Request.Context(), discID, userID, dto.Content)
    if err != nil {
        logger.Errorf("failed to add comment: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add comment"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"id": commentID})
}

// GET /discussions/:id/comments
func (ctr *Controller) List(c *gin.Context) {
    discID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discussion ID"})
        return
    }

    comments, err := ctr.svc.GetComments(c.Request.Context(), discID)
    if err != nil {
        logger.Errorf("failed to list comments: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch comments"})
        return
    }

    c.JSON(http.StatusOK, comments)
}
