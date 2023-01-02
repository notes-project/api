package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// health endpoints for k8s probes
const (
	healthPort = "3040"

	readinessEndpoint = "/readyz"
	livenessEndpoint  = "/healthz"
)

func (s server) isReady(c *gin.Context) {
	isReady := s.db.IsReady()

	if !isReady {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}

}
func (s server) isAlive(c *gin.Context) {
	c.Status(http.StatusOK)
}
