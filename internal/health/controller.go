// controller.go 
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	service *HealthService
}

func NewHealthController(service *HealthService) *HealthController {
	return &HealthController{service: service}
}

func (hc *HealthController) HandleHealthCheck(c *gin.Context) {
	status := hc.service.CheckHealth()

	if status.Status == "ok" {
		c.JSON(http.StatusOK, status)
	} else {
		c.JSON(http.StatusServiceUnavailable, status)
	}
}
