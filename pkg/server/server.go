package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server interface {
	Start() error
}

type server struct {
	serverConfiguration
	logger *zap.Logger

	servers           *[]*http.Server
	serversShutdowned chan bool
	serverError       chan error
}

func (s server) Start() error {

	s.servers = new([]*http.Server)
	s.serversShutdowned = make(chan bool, 1)
	s.serverError = make(chan error, 1)

	s.startMainServers()
	s.serveHealthProbes()

	go s.handleGracefulShutdown()

	select {
	case <-s.serversShutdowned:
		s.logger.Info("Servers shut down successfully")
		return nil
	case err := <-s.serverError:
		return err
	}

}

func (s server) handleGracefulShutdown() {
	quit := make(chan os.Signal, 2)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, server := range *s.servers {
		err := server.Shutdown(ctx)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Server on port '%s' failed during shutdown, error: %s", server.Addr, err))
		}
	}

	s.serversShutdowned <- true
}

func (s server) startMainServers() {
	defaultRouter := gin.Default()
	defaultRouter.SetTrustedProxies(nil)

	v1 := defaultRouter.Group("/api/v1")
	{
		v1.POST("/notes", s.addNote)
		v1.GET("/notes", s.getNotes)

		v1.GET("/notes/:title", s.getNoteByTitle)
		v1.DELETE("/notes/:title", s.deleteNoteByTitle)
		v1.POST("/notes/:title", s.updateNoteByTitle)
	}

	s.serveHttp(defaultRouter)
	s.serveHttps(defaultRouter)
}

func (s server) serveHttp(router *gin.Engine) {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: router,
	}

	*s.servers = append(*s.servers, httpServer)

	go func() {
		s.logger.Info(fmt.Sprintf("Server started on port %s for HTTP", s.port))

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(fmt.Sprintf("HTTP server crashed, error: %s", err))
			s.serverError <- err
		}
	}()
}

func (s server) serveHttps(router *gin.Engine) {
	httpsServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.tlsPort),
		Handler: router,
	}

	if len(s.tlsPort) > 0 && s.tlsCertLocation != "" && s.tlsKeyLocation != "" {

		*s.servers = append(*s.servers, httpsServer)

		go func() {
			s.logger.Info(fmt.Sprintf("Server started on port %s for HTTPS", s.tlsPort))

			err := httpsServer.ListenAndServeTLS(s.tlsCertLocation, s.tlsKeyLocation)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.logger.Error(fmt.Sprintf("HTTPS server crashed, error: %s", err))
				s.serverError <- err
			}
		}()
	}
}
