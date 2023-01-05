package database

import (
	"context"
	"fmt"

	"github.com/notes-project/api/pkg/adapters"
	facademongo "github.com/notes-project/api/pkg/facade/go.mongodb.org/mongo-driver/mongo"
	"github.com/notes-project/api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type Database interface {
	Connect() error

	IsReady() bool

	AddNote(note model.Note) error
	UpdateNote(noteTitle string, updatedNote model.Note) error
	GetNote(noteTitle string) (model.Note, error)
	GetNotes() ([]model.Note, error)
	GetNotesFiltered(tags []string, category, date string) ([]model.Note, error)
	DeleteNote(noteTitle string) error
	DeleteNotes() error
}

type database struct {
	databaseConfiguration
	logger *zap.Logger

	// populated automatically inside the Connect() method
	client     adapters.DbClient
	collection adapters.DbCollection
}

const (
	// used as a primary key
	noteTitlePrimaryKey = "title"
)

var (
	ctx = context.Background()
)

func (d *database) Connect() error {

	if d.client != nil {
		d.logger.Info("Database already started")
		return nil
	}

	var err error

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(d.connectionUri).SetServerAPIOptions(serverAPIOptions)

	d.client, err = facademongo.GetClientInstace().Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to the database, error: %w", err)
	}

	err = facademongo.GetClientInstace().Ping(d.client.(*mongo.Client), ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to verify database connection, error: %w", err)
	}

	db := facademongo.GetClientInstace().Database(d.client.(*mongo.Client), d.databaseName)

	d.collection = facademongo.GetDatabaseInstace().Collection(db, d.collectionName)

	err = d.setUniqueIndexes()
	if err != nil {
		return err
	}

	d.logger.Info("Successfully connected to the database")

	return nil
}

func (d *database) setUniqueIndexes() error {
	indexView := d.collection.Indexes()

	_, err := facademongo.GetIndexViewInstace().CreateOne(indexView, ctx, mongo.IndexModel{
		// the title filed of the note from model.Note
		Keys: bson.D{
			{Key: noteTitlePrimaryKey, Value: -1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to set '%s' as a unique collection index, error: %w", noteTitlePrimaryKey, err)
	}

	d.logger.Info(fmt.Sprintf("Successfully set '%s' as a unique collection index", noteTitlePrimaryKey))

	return nil
}

func (d *database) IsReady() bool {
	err := d.client.Ping(ctx, readpref.Primary())
	return err == nil
}
