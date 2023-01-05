package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/notes-project/api/pkg/constants"
	"github.com/notes-project/api/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

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
			s.logger.Info(fmt.Sprintf("Note '%s' already exists in database", note.Title))

			c.JSON(http.StatusBadRequest,
				gin.H{
					"error": fmt.Sprintf("note with key 'title' and value '%s' already exists", note.Title),
				},
			)

			return
		}

		s.logger.Error(fmt.Sprintf("Failed to add a note to the database, err: %s", err))

		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}
}

func (s server) getNotes(c *gin.Context) {
	var notes []model.Note
	var err error

	tags := strings.Split(c.Query("tags"), ",")
	category := c.Query("category")
	date := c.Query("date")

	notes, err = s.db.GetNotesFiltered(tags, category, date)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get notes from database, err: %s", err))

		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": "failed to retrieve notes",
			},
		)

		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"notes": notes,
		})
}

func (s server) getNoteByTitle(c *gin.Context) {
	noteTtile := c.Param("title")

	note, err := s.db.GetNote(noteTtile)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Info(fmt.Sprintf("Note '%s' does not exist in database", noteTtile))

			c.JSON(http.StatusNotFound,
				gin.H{
					"error": fmt.Sprintf("note '%s' does not exist", noteTtile),
				},
			)

			return
		}

		s.logger.Error(fmt.Sprintf("Failed to get note '%s' from database, err: %s", noteTtile, err))

		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("failed to retrieve note '%s'", noteTtile),
			},
		)

		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"note": note,
		})
}

func (s server) deleteNoteByTitle(c *gin.Context) {
	noteTtile := c.Param("title")

	err := s.db.DeleteNote(noteTtile)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Info(fmt.Sprintf("Note '%s' does not exist in database", noteTtile))

			c.JSON(http.StatusNoContent,
				gin.H{
					"info": fmt.Sprintf("note '%s' does not exist", noteTtile),
				},
			)

			return
		}

		s.logger.Error(fmt.Sprintf("Failed to delete note '%s' from database, err: %s", noteTtile, err))

		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("failed to retrieve note '%s'", noteTtile),
			},
		)

		return
	}

	c.Status(http.StatusOK)
}

func (s server) deleteNotes(c *gin.Context) {

	err := s.db.DeleteNotes()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Info("0 notes in database")

			c.JSON(http.StatusNoContent,
				gin.H{
					"info": "no notes available",
				},
			)

			return
		}

		s.logger.Error(fmt.Sprintf("Failed to delete notes from database, err: %s", err))

		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": "failed to delete notes",
			},
		)

		return
	}

	c.Status(http.StatusOK)
}

func (s server) updateNoteByTitle(c *gin.Context) {
	noteTtile := c.Param("title")

	note := model.Note{}

	err := c.MustBindWith(&note, binding.JSON)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	note.Date = time.Now().Format(constants.DateFormat)

	err = s.db.UpdateNote(noteTtile, note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Info(fmt.Sprintf("Note '%s' does not exist in database", noteTtile))

			c.JSON(http.StatusNotFound,
				gin.H{
					"error": fmt.Sprintf("note '%s' does not exist", noteTtile),
				},
			)

			return
		}

		s.logger.Error(fmt.Sprintf("Failed to update note '%s' from database, err: %s", noteTtile, err))

		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("failed to update note '%s'", noteTtile),
			},
		)

		return
	}

	c.Status(http.StatusOK)
}
