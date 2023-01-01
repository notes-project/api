package database

import (
	"fmt"

	"github.com/denislavPetkov/notes/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (d *database) AddNote(note model.Note) error {
	_, err := d.collection.InsertOne(ctx, note)

	if err != nil {
		return fmt.Errorf("failed to add note %v to the collection, error: %w", note, err)
	}

	d.logger.Info("Successfully added note to the collection", zap.Any("note", note))

	return nil
}

func (d *database) UpdateNote(noteTitle string, updatedNote model.Note) error {
	result, err := d.collection.ReplaceOne(ctx, bson.D{
		{
			Key:   noteTitlePrimaryKey,
			Value: noteTitle,
		},
	},
		updatedNote,
	)
	if err != nil {
		return fmt.Errorf("failed to update note '%s', error: %w", noteTitle, err)
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (d *database) GetNote(noteTitle string) (model.Note, error) {
	result := d.collection.FindOne(ctx, bson.D{
		{
			Key:   noteTitlePrimaryKey,
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

func (d *database) GetNotesFiltered(tags []string, category, date string) ([]model.Note, error) {

	tagsFilter := getTagsFilter(tags)
	categoryFilter := getCategoryFilter(category)
	dateFilter := getDateFilter(date)

	cursor, err := d.collection.Find(ctx, bson.D{
		tagsFilter,
		categoryFilter,
		dateFilter,
	})
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

func getTagsFilter(tags []string) bson.E {
	if len(tags) == 1 && tags[0] == "" {
		return bson.E{}
	}

	return bson.E{
		Key: "tags",
		Value: bson.D{
			{
				Key:   "$all",
				Value: tags,
			},
		},
	}
}

func getDateFilter(date string) bson.E {
	if date == "" {
		return bson.E{}
	}

	return bson.E{
		Key:   "date",
		Value: date,
	}
}

func getCategoryFilter(category string) bson.E {
	if category == "" {
		return bson.E{}
	}

	return bson.E{
		Key:   "category",
		Value: category,
	}
}

func (d *database) DeleteNote(noteTitle string) error {
	result, err := d.collection.DeleteOne(ctx,
		bson.D{
			{
				Key:   noteTitlePrimaryKey,
				Value: noteTitle,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to delete note '%s' from collection, error: %w", noteTitle, err)
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
