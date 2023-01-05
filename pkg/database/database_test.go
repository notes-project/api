package database

import (
	"errors"

	"github.com/golang/mock/gomock"
	facademongo "github.com/notes-project/api/pkg/facade/go.mongodb.org/mongo-driver/mongo"
	mockadapters "github.com/notes-project/api/pkg/mock/adapters"
	mockfacademongo "github.com/notes-project/api/pkg/mock/facade/go.mongodb.org/mongo-driver/mongo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var _ = Describe("Database", func() {

	var (
		ctrl *gomock.Controller

		mockFacadeIndexView   *mockfacademongo.MockIndexView
		mockFacadeDatabase    *mockfacademongo.MockDatabase
		mockFacadeMongoClient *mockfacademongo.MockClient
		mockDbClient          *mockadapters.MockDbClient
		mockDbCollection      *mockadapters.MockDbCollection

		dbInstance *database
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockFacadeIndexView = mockfacademongo.NewMockIndexView(ctrl)
		mockFacadeDatabase = mockfacademongo.NewMockDatabase(ctrl)
		mockFacadeMongoClient = mockfacademongo.NewMockClient(ctrl)

		mockDbClient = mockadapters.NewMockDbClient(ctrl)
		mockDbCollection = mockadapters.NewMockDbCollection(ctrl)

		dbInstance = &database{
			databaseConfiguration: databaseConfiguration{
				connectionUri:  "",
				databaseName:   "",
				collectionName: "",
			},
			logger:     zap.L(),
			collection: mockDbCollection,
		}

		facademongo.SetIndexViewInstance(mockFacadeIndexView)
		facademongo.SetClientInstance(mockFacadeMongoClient)
		facademongo.SetDatabaseInstance(mockFacadeDatabase)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Connect", func() {
		It("should return nil when client already initialized", func() {
			dbInstance.client = mockDbClient

			err := dbInstance.Connect()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error when failed to connect", func() {
			mockFacadeMongoClient.EXPECT().Connect(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

			err := dbInstance.Connect()
			Expect(err).To(HaveOccurred())
		})

		It("should return an error when failed to ping database", func() {
			mockFacadeMongoClient.EXPECT().Connect(gomock.Any(), gomock.Any()).Return(&mongo.Client{}, nil)
			mockFacadeMongoClient.EXPECT().Ping(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))

			err := dbInstance.Connect()
			Expect(err).To(HaveOccurred())
		})

		It("should return an error when failed to create a new indexe", func() {
			mockFacadeMongoClient.EXPECT().Connect(gomock.Any(), gomock.Any()).Return(&mongo.Client{}, nil)
			mockFacadeMongoClient.EXPECT().Ping(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockFacadeMongoClient.EXPECT().Database(gomock.Any(), gomock.Any()).Return(&mongo.Database{})
			mockFacadeDatabase.EXPECT().Collection(gomock.Any(), gomock.Any()).Return(&mongo.Collection{})
			mockFacadeIndexView.EXPECT().CreateOne(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New(""))

			err := dbInstance.Connect()
			Expect(err).To(HaveOccurred())
		})

		It("should return nil when no error occurred", func() {
			mockFacadeMongoClient.EXPECT().Connect(gomock.Any(), gomock.Any()).Return(&mongo.Client{}, nil)
			mockFacadeMongoClient.EXPECT().Ping(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockFacadeMongoClient.EXPECT().Database(gomock.Any(), gomock.Any()).Return(&mongo.Database{})
			mockFacadeDatabase.EXPECT().Collection(gomock.Any(), gomock.Any()).Return(&mongo.Collection{})
			mockFacadeIndexView.EXPECT().CreateOne(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)

			err := dbInstance.Connect()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("IsReady", func() {
		BeforeEach(func() {
			dbInstance.client = mockDbClient
		})

		It("should return false when ping returns an error", func() {
			mockDbClient.EXPECT().Ping(gomock.Any(), gomock.Any()).Return(errors.New(""))

			isReady := dbInstance.IsReady()
			Expect(isReady).To(BeFalse())
		})

		It("should return true when ping returns no error", func() {
			mockDbClient.EXPECT().Ping(gomock.Any(), gomock.Any()).Return(nil)

			isReady := dbInstance.IsReady()
			Expect(isReady).To(BeTrue())
		})
	})

})
