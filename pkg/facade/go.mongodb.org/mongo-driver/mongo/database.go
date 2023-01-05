package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Collection(db *mongo.Database, name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type database struct{}

var databaseInstance Database = database{}

func SetDatabaseInstance(c Database) {
	databaseInstance = c
}

func GetDatabaseInstace() Database {
	return databaseInstance
}

func (d database) Collection(db *mongo.Database, name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return db.Collection(name, opts...)
}
