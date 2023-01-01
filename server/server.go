package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/denislavPetkov/notes/constants"
	"github.com/denislavPetkov/notes/database"
	"github.com/denislavPetkov/notes/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server interface {
	Start()
}

type server struct {
	port      string
	ginRouter *gin.Engine
	db        database.Database
	logger    *zap.Logger
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
	}

	s.logger.Info(fmt.Sprintf("Server started on port %s", s.port))

	s.ginRouter.Run(fmt.Sprintf(":%s", s.port))
}

func (s server) addNote(c *gin.Context) {

	note := model.Note{}

	err := c.MustBindWith(&note, binding.JSON)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	note.Date = time.Now().Format(constants.DateFormat)

	err = s.db.AddNote(note)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest,
				gin.H{
					"error": fmt.Sprintf("note with key 'title' and value '%s' already exists", note.Title),
				},
			)

			return
		}

		s.logger.Error(fmt.Sprintf("Failed to add a note to the database, err: %s", err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

}

func (s server) getNotes(c *gin.Context) {
	notes, err := s.db.GetNotes()
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get notes from database, err: %s", err))

		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": "failed to retrieve notes from database",
			},
		)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})

}
