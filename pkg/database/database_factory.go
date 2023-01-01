package database

import (
	"sync"

	"go.uber.org/zap"
)

type DatabaseFactory interface {
	NewDatabase(connectionUri, dbName, dbCollection string) Database
}

type databaseFactory struct{}

func NewDatabaseFactory() DatabaseFactory {
	return databaseFactory{}
}

var (
	once sync.Once

	databaseInstance Database
)

func (df databaseFactory) NewDatabase(connectionUri, dbName, dbCollection string) Database {

	once.Do(func() {
		databaseInstance = &database{
			connectionUri:  connectionUri,
			databaseName:   dbName,
			collectionName: dbCollection,
			logger:         zap.L().Named("Database"),
		}
	})

	return databaseInstance
}
