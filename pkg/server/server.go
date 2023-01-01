package server

import (
	"fmt"

	"github.com/denislavPetkov/notes/pkg/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server interface {
	Start()
}

type server struct {
	port   string
	db     database.Database
	logger *zap.Logger

	// populated automatically by the Start() method
	ginRouter *gin.Engine
}

func (s server) Start() {
	s.ginRouter = gin.Default()

	base := s.ginRouter.Group("/api")
	{
		base.GET(readinessEndpoint, s.isReady)
		base.GET(livenessEndpoint, s.isAlive)
	}

	v1 := s.ginRouter.Group("/api/v1")
	{
		v1.POST("/notes", s.addNote)
		v1.GET("/notes", s.getNotes)

		v1.GET("/notes/:title", s.getNoteByTitle)
		v1.DELETE("/notes/:title", s.deleteNoteByTitle)
		v1.POST("/notes/:title", s.updateNoteByTitle)
	}

	s.logger.Info(fmt.Sprintf("Server started on port %s", s.port))

	s.ginRouter.Run(fmt.Sprintf(":%s", s.port))
}
