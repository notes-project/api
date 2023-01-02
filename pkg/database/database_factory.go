package database

import (
	"sync"

	"go.uber.org/zap"
)

type DatabaseFactory interface {
	NewDatabase(dbConfig databaseConfiguration) Database
}

type databaseFactory struct{}

func NewDatabaseFactory() DatabaseFactory {
	return databaseFactory{}
}

var (
	once sync.Once

	databaseInstance Database
)

func (df databaseFactory) NewDatabase(dbConfig databaseConfiguration) Database {

	once.Do(func() {
		databaseInstance = &database{
			databaseConfiguration: dbConfig,
			logger:                zap.L().Named("Database"),
		}
	})

	return databaseInstance
}
