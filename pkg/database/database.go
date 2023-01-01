package database

import (
	"context"
	"fmt"

	"github.com/denislavPetkov/notes/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type Database interface {
	Start() error

	IsReady() bool

	AddNote(note model.Note) error
	UpdateNote() error
	GetNote() error
	GetNotes() ([]model.Note, error)
	DeleteNote() error
}

type database struct {
	connectionUri string
	databaseName  string
	collection    string

	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection

	logger *zap.Logger
}

var (
	ctx = context.TODO()
)

func (d *database) Start() error {

	if d.mongoClient != nil {
		d.logger.Info("Database already started")
		return nil
	}

	var err error
	// set auth?
	clientOptions := options.Client().ApplyURI(d.connectionUri)

	d.mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to the database, err: %w", err)
	}

	err = d.mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to verify database connection, err: %w", err)
	}

	d.mongoCollection = d.mongoClient.Database(d.databaseName).Collection(d.collection)

	err = d.setUniqueIndexes()
	if err != nil {
		return err
	}

	d.logger.Info("Successfully connected to the database")

	return nil
}

func (d *database) setUniqueIndexes() error {
	noteTitleKey := "title"

	_, err := d.mongoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		// the title filed of the note from model.Note
		Keys: bson.D{
			{Key: noteTitleKey, Value: -1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to set '%s' as a unique collection index, err: %w", noteTitleKey, err)
	}

	d.logger.Info(fmt.Sprintf("Successfully set '%s' as a unique collection index", noteTitleKey))

	return nil
}

func (d *database) IsReady() bool {
	err := d.mongoClient.Ping(context.TODO(), readpref.Primary())
	return err == nil
}

func (d *database) AddNote(note model.Note) error {
	_, err := d.mongoCollection.InsertOne(ctx, note)

	if err != nil {
		return fmt.Errorf("failed to add note %v to the collection, err: %w", note, err)
	}

	d.logger.Info("Successfully added note to the collection", zap.Any("note", note))

	return nil
}

func (d *database) UpdateNote() error { return nil }
func (d *database) GetNote() error    { return nil }

func (d *database) GetNotes() ([]model.Note, error) {
	cursor, err := d.mongoCollection.Find(ctx, bson.D{})
	if err != nil {
		return []model.Note{}, fmt.Errorf("failed to get notes from collection, err: %w", err)
	}

	var notes []model.Note
	err = cursor.All(ctx, &notes)
	if err != nil {
		return []model.Note{}, fmt.Errorf("failed to get notes from collection, err: %w", err)
	}

	return notes, nil
}

func (d *database) DeleteNote() error { return nil }
