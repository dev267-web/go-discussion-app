// controller.go 
package subscription

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-discussion-app/models"
)

type SubscriptionController struct {
	service *Service
}

func NewSubscriptionController(service *Service) *SubscriptionController {
	return &SubscriptionController{service}
}

// POST /discussions/:id/subscribe
func (sc *SubscriptionController) Subscribe(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discussion ID"})
		return
	}

	var subDTO SubscribeDTO
	if err := c.ShouldBindJSON(&subDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId") // Optional
	uid, _ := userID.(int)

	sub := &models.Subscription{
		DiscussionID: discussionID,
		Email:        subDTO.Email,
		SubscribedAt: subDTO.SubscribedAt,
	}
	if uid != 0 {
		sub.UserID = &uid
	}

	if err := sc.service.Subscribe(sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "subscribed successfully"})
}

// DELETE /discussions/:id/unsubscribe
func (sc *SubscriptionController) Unsubscribe(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discussion ID"})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sc.service.Unsubscribe(discussionID, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unsubscribed successfully"})
}

// POST /discussions/:id/notify
func (sc *SubscriptionController) Notify(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discussion ID"})
		return
	}

	var req struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sc.service.NotifySubscribers(discussionID, req.Subject, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notifications sent"})
}
