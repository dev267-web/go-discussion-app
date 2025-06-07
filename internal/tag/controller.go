// controller.go 
package tag

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "go-discussion-app/pkg/logger"
)

// TagController handles HTTP requests for tags.
type TagController struct {
    svc *TagService
}

// NewController constructs a TagController.
func NewController(svc *TagService) *TagController {
    return &TagController{svc: svc}
}

// ListHandler handles GET /tags
func (ctr *TagController) ListHandler(c *gin.Context) {
    tags, err := ctr.svc.ListTags(c.Request.Context())
    if err != nil {
        logger.Errorf("failed to list tags: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
        return
    }
    c.JSON(http.StatusOK, tags)
}
