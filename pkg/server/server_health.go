package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// health endpoints for k8s probes
const (
	healthPort = "3040"

	readinessEndpoint = "/readyz"
	livenessEndpoint  = "/healthz"
)

func (s server) serveHealthProbes() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET(readinessEndpoint, s.isReady)
	router.GET(livenessEndpoint, s.isAlive)

	healthServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", healthPort),
		Handler: router,
	}

	*s.servers = append(*s.servers, healthServer)

	s.logger.Info(fmt.Sprintf("Serving probes on port %s", healthPort))

	go func() {
		err := healthServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(fmt.Sprintf("HTTP Health server crashed, error: %s", err))
			s.serverError <- err
		}
	}()
}

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
