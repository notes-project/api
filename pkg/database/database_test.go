package database

import (
	"errors"

	facademongo "github.com/denislavPetkov/notes/pkg/facade/go.mongodb.org/mongo-driver/mongo"
	mockdatabase "github.com/denislavPetkov/notes/pkg/mock/database"
	mockmongo "github.com/denislavPetkov/notes/pkg/mock/facade/go.mongodb.org/mongo-driver/mongo"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Database", func() {

	var (
		ctrl *gomock.Controller

		mockMongoClient  *mockmongo.MockClient
		mockDbClient     *mockdatabase.MockDbClient
		mockDbCollection *mockdatabase.MockDbCollection

		dbInstance *database
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockMongoClient = mockmongo.NewMockClient(ctrl)
		mockDbClient = mockdatabase.NewMockDbClient(ctrl)
		mockDbCollection = mockdatabase.NewMockDbCollection(ctrl)

		dbInstance = &database{
			databaseConfiguration: databaseConfiguration{
				connectionUri:  "",
				databaseName:   "",
				collectionName: "",
			},
			logger:     zap.L(),
			collection: mockDbCollection,
		}

		facademongo.SetInstance(mockMongoClient)
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
			mockMongoClient.EXPECT().Connect(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

			err := dbInstance.Connect()
			Expect(err).To(HaveOccurred())
		})
	})

})
