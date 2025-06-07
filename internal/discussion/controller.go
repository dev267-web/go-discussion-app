// controller.go 
package discussion

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

// POST /discussions
func (ctr *Controller) Create(c *gin.Context) {
    userID, _ := auth.GetUserID(c)
    var dto CreateDiscussionDTO
    if err := c.ShouldBindJSON(&dto); err != nil || dto.Validate() != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    id, err := ctr.svc.Create(c.Request.Context(), userID, &dto)
    if err != nil {
        logger.Errorf("create discussion error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GET /discussions
func (ctr *Controller) List(c *gin.Context) {
    ds, err := ctr.svc.GetAll(c.Request.Context())
    if err != nil {
        logger.Errorf("list discussions error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list"})
        return
    }
    c.JSON(http.StatusOK, ds)
}

// GET /discussions/:id
func (ctr *Controller) Get(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    d, err := ctr.svc.GetByID(c.Request.Context(), id)
    if err != nil {
        logger.Errorf("get discussion error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch"})
        return
    }
    if d == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, d)
}

// PUT /discussions/:id
func (ctr *Controller) Update(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var dto UpdateDiscussionDTO
    if err := c.ShouldBindJSON(&dto); err != nil || dto.Validate() != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    d, err := ctr.svc.Update(c.Request.Context(), id, &dto)
    if err != nil {
        logger.Errorf("update discussion error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update"})
        return
    }
    c.JSON(http.StatusOK, d)
}

// DELETE /discussions/:id
func (ctr *Controller) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := ctr.svc.Delete(c.Request.Context(), id); err != nil {
        logger.Errorf("delete discussion error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete"})
        return
    }
    c.Status(http.StatusNoContent)
}

// GET /discussions/user/:userId
func (ctr *Controller) ListByUser(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("userId"))
    ds, err := ctr.svc.GetByUser(c.Request.Context(), uid)
    if err != nil {
        logger.Errorf("list by user error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list"})
        return
    }
    c.JSON(http.StatusOK, ds)
}

// GET /discussions/tag/:tag
func (ctr *Controller) ListByTag(c *gin.Context) {
    tag := c.Param("tag")
    ds, err := ctr.svc.GetByTag(c.Request.Context(), tag)
    if err != nil {
        logger.Errorf("list by tag error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list"})
        return
    }
    c.JSON(http.StatusOK, ds)
}

// POST /discussions/:id/tags
func (ctr *Controller) AddTags(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var dto AddTagsDTO
    if err := c.ShouldBindJSON(&dto); err != nil || dto.Validate() != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    if err := ctr.svc.AddTags(c.Request.Context(), id, &dto); err != nil {
        logger.Errorf("add tags error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add tags"})
        return
    }
    c.Status(http.StatusNoContent)
}

// POST /discussions/schedule
func (ctr *Controller) Schedule(c *gin.Context) {
    userID, _ := auth.GetUserID(c)
    var dto ScheduleDTO
    if err := c.ShouldBindJSON(&dto); err != nil || dto.Validate() != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    id, err := ctr.svc.Schedule(c.Request.Context(), userID, &dto)
    if err != nil {
        logger.Errorf("schedule discussion error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not schedule"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"id": id})
}
