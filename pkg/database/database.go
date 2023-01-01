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
	GetNote(noteTitle string) (model.Note, error)
	GetNotes() ([]model.Note, error)
	DeleteNote(noteTitle string) error
}

type database struct {
	connectionUri  string
	databaseName   string
	collectionName string

	client     *mongo.Client
	collection *mongo.Collection

	logger *zap.Logger
}

const (
	noteTitleKey = "title"
)

var (
	ctx = context.TODO()
)

func (d *database) Start() error {

	if d.client != nil {
		d.logger.Info("Database already started")
		return nil
	}

	var err error
	// set auth?
	clientOptions := options.Client().ApplyURI(d.connectionUri)

	d.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to the database, error: %w", err)
	}

	err = d.client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to verify database connection, error: %w", err)
	}

	d.collection = d.client.Database(d.databaseName).Collection(d.collectionName)

	err = d.setUniqueIndexes()
	if err != nil {
		return err
	}

	d.logger.Info("Successfully connected to the database")

	return nil
}

func (d *database) setUniqueIndexes() error {
	_, err := d.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		// the title filed of the note from model.Note
		Keys: bson.D{
			{Key: noteTitleKey, Value: -1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to set '%s' as a unique collection index, error: %w", noteTitleKey, err)
	}

	d.logger.Info(fmt.Sprintf("Successfully set '%s' as a unique collection index", noteTitleKey))

	return nil
}

func (d *database) IsReady() bool {
	err := d.client.Ping(context.TODO(), readpref.Primary())
	return err == nil
}

func (d *database) AddNote(note model.Note) error {
	_, err := d.collection.InsertOne(ctx, note)

	if err != nil {
		return fmt.Errorf("failed to add note %v to the collection, error: %w", note, err)
	}

	d.logger.Info("Successfully added note to the collection", zap.Any("note", note))

	return nil
}

func (d *database) UpdateNote() error { return nil }

func (d *database) GetNote(noteTitle string) (model.Note, error) {

	result := d.collection.FindOne(ctx, bson.D{
		{
			Key:   noteTitleKey,
			Value: noteTitle,
		},
	},
	)

	note := model.Note{}

	err := result.Decode(&note)
	if err != nil {
		return model.Note{}, fmt.Errorf("failed to decode note into object, error: %w", err)
	}

	return note, nil
}

func (d *database) GetNotes() ([]model.Note, error) {

	cursor, err := d.collection.Find(ctx, bson.D{})
	if err != nil {
		return []model.Note{}, fmt.Errorf("failed to get notes from collection, error: %w", err)
	}

	var notes []model.Note
	err = cursor.All(ctx, &notes)
	if err != nil {
		return []model.Note{}, fmt.Errorf("failed to get notes from collection, error: %w", err)
	}

	return notes, nil
}

func (d *database) DeleteNote(noteTitle string) error {

	_, err := d.collection.DeleteOne(ctx,
		bson.D{
			{
				Key:   noteTitleKey,
				Value: noteTitle,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to delete note '%s' from collection, error: %w", noteTitle, err)
	}

	return nil
}
