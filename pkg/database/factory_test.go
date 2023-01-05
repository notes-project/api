package database

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DatabaseFactory", func() {

	var (
		databaseFactory DatabaseFactory = databaseFactory{}
	)

	Describe("NewDatabaseFactory", func() {
		It("should return a new database factory object", func() {
			dbFactory := NewDatabaseFactory()
			Expect(dbFactory).NotTo(BeNil())
		})
	})

	Describe("NewDatabase", func() {
		It("should return a new database object", func() {
			dbInstance := databaseFactory.NewDatabase(databaseConfiguration{})
			Expect(dbInstance).NotTo(BeNil())
		})
	})

})
