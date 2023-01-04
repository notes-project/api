package database

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DatabaseConfig", func() {

	Describe("NewDatabaseConfiguration", func() {
		connectionUri := "connectionUri"
		databaseName := "databaseName"
		collectionName := "collectionName"

		It("should return a new databese configuration object", func() {
			dbConfig := NewDatabaseConfiguration(connectionUri, databaseName, collectionName)

			Expect(dbConfig).NotTo(BeNil())
			Expect(dbConfig.connectionUri).To(Equal(connectionUri))
			Expect(dbConfig.databaseName).To(Equal(databaseName))
			Expect(dbConfig.collectionName).To(Equal(collectionName))
		})
	})

})
