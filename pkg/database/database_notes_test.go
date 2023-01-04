package database

import (
	"errors"

	mockdatabase "github.com/denislavPetkov/notes/pkg/mock/database"
	"github.com/denislavPetkov/notes/pkg/model"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var _ = Describe("DatabaseNotes", func() {

	var (
		ctrl *gomock.Controller

		mockDbClient     *mockdatabase.MockDbClient
		mockDbCollection *mockdatabase.MockDbCollection

		dbInstance *database
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockDbClient = mockdatabase.NewMockDbClient(ctrl)
		mockDbCollection = mockdatabase.NewMockDbCollection(ctrl)

		dbInstance = &database{
			databaseConfiguration: databaseConfiguration{
				connectionUri:  "",
				databaseName:   "",
				collectionName: "",
			},
			logger:     zap.L(),
			client:     mockDbClient,
			collection: mockDbCollection,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("AddNote", func() {
		It("should add a note to the database when error does not occur", func() {
			mockDbCollection.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(nil, nil)

			err := dbInstance.AddNote(model.Note{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error when failed to insert to database", func() {
			mockDbCollection.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

			err := dbInstance.AddNote(model.Note{})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("UpdateNote", func() {
		It("should update a note from the database when error does not occur and note is in database", func() {
			mockDbCollection.EXPECT().ReplaceOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{
				MatchedCount: 1,
			}, nil)

			err := dbInstance.UpdateNote("", model.Note{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error when failed to update note in database", func() {
			mockDbCollection.EXPECT().ReplaceOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

			err := dbInstance.UpdateNote("", model.Note{})
			Expect(err).To(HaveOccurred())
		})

		It("should return an error when note to update is not in database", func() {
			mockDbCollection.EXPECT().ReplaceOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{
				MatchedCount: 0,
			}, nil)

			err := dbInstance.UpdateNote("", model.Note{})
			Expect(err).To(MatchError(mongo.ErrNoDocuments))
		})
	})

	Describe("GetNote", func() {
		It("should return note when no error occurs", func() {
			mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				mongo.NewSingleResultFromDocument(model.Note{
					Title: "test",
				}, nil, nil),
			)

			note, err := dbInstance.GetNote("test")

			Expect(err).NotTo(HaveOccurred())
			Expect(note).NotTo(BeNil())
			Expect(note.Title).To(Equal("test"))
		})

		It("should return an error when failed to get note", func() {
			mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				mongo.NewSingleResultFromDocument(nil, errors.New(""), nil),
			)

			_, err := dbInstance.GetNote("")

			Expect(err).To(HaveOccurred())
		})

		It("should return an error when failed to get note", func() {
			mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				mongo.NewSingleResultFromDocument(nil, errors.New(""), nil),
			)

			_, err := dbInstance.GetNote("")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetNotes", func() {
		It("should return notes when no error occurs", func() {
			mockDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mongo.NewCursorFromDocuments(
				[]interface{}{
					model.Note{
						Title: "test1",
					},
					model.Note{
						Title: "test2",
					},
				},
				nil, nil),
			)

			notes, err := dbInstance.GetNotes()

			Expect(err).NotTo(HaveOccurred())
			Expect(notes).NotTo(BeEmpty())
			Expect(notes).To(HaveLen(2))
		})

		It("should return an error when failed to get notes", func() {
			mockDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mongo.NewCursorFromDocuments([]interface{}{nil}, nil, nil))

			_, err := dbInstance.GetNotesFiltered(nil, "", "")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetNotesFiltered", func() {
		It("should return notes when no error occurs", func() {
			mockDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mongo.NewCursorFromDocuments(
				[]interface{}{
					model.Note{
						Title: "test1",
					},
					model.Note{
						Title: "test2",
					},
				},
				nil, nil),
			)

			notes, err := dbInstance.GetNotesFiltered(nil, "", "")

			Expect(err).NotTo(HaveOccurred())
			Expect(notes).NotTo(BeEmpty())
			Expect(notes).To(HaveLen(2))
		})

		It("should return an error when failed to get notes", func() {
			mockDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mongo.NewCursorFromDocuments([]interface{}{nil}, nil, nil))

			_, err := dbInstance.GetNotes()

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteNote", func() {
		It("should return no error when no error occurs", func() {
			mockDbCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any()).Return(
				&mongo.DeleteResult{
					DeletedCount: 1,
				},
				nil,
			)

			err := dbInstance.DeleteNote("")

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when failed to delete note", func() {
			mockDbCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any()).Return(
				&mongo.DeleteResult{},
				errors.New(""),
			)

			err := dbInstance.DeleteNote("")

			Expect(err).To(HaveOccurred())
		})

		It("should return error when no notes in database", func() {
			mockDbCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any()).Return(
				&mongo.DeleteResult{
					DeletedCount: 0,
				},
				nil,
			)

			err := dbInstance.DeleteNote("")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteNotes", func() {
		It("should return no error when no error occurs", func() {
			mockDbCollection.EXPECT().DeleteMany(gomock.Any(), gomock.Any()).Return(
				&mongo.DeleteResult{
					DeletedCount: 1,
				},
				nil,
			)

			err := dbInstance.DeleteNotes()

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when failed to delete note", func() {
			mockDbCollection.EXPECT().DeleteMany(gomock.Any(), gomock.Any()).Return(
				&mongo.DeleteResult{},
				errors.New(""),
			)

			err := dbInstance.DeleteNotes()

			Expect(err).To(HaveOccurred())
		})

		It("should return error when no notes in database", func() {
			mockDbCollection.EXPECT().DeleteMany(gomock.Any(), gomock.Any()).Return(
				&mongo.DeleteResult{
					DeletedCount: 0,
				},
				nil,
			)

			err := dbInstance.DeleteNotes()

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("getTagsFilter", func() {
		It("should return an empty object when no tags provided", func() {
			filter := getTagsFilter([]string{""})
			Expect(filter).To(Equal(bson.E{}))
		})

		It("should return not empty object when tags are provided", func() {
			filter := getTagsFilter([]string{"test"})

			Expect(filter).NotTo(Equal(bson.E{}))
			Expect(filter).To(Equal(bson.E{
				Key: "tags",
				Value: bson.D{
					{
						Key:   "$all",
						Value: []string{"test"},
					},
				},
			}))
		})
	})

	Describe("getDateFilter", func() {
		It("should return empty object when no date provided", func() {
			filter := getDateFilter("")
			Expect(filter).To(Equal(bson.E{}))
		})

		It("should return not empty object when date is provided", func() {
			filter := getDateFilter("test")
			Expect(filter).To(Equal(bson.E{
				Key:   "date",
				Value: "test",
			}))
		})
	})

	Describe("getCategoryFilter", func() {
		It("should return empty object when no category provided", func() {
			filter := getCategoryFilter("")
			Expect(filter).To(Equal(bson.E{}))
		})

		It("should return not empty object when tags category is provided", func() {
			filter := getCategoryFilter("test")
			Expect(filter).To(Equal(bson.E{
				Key:   "category",
				Value: "test",
			}))
		})
	})

})
