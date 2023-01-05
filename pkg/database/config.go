package database

type databaseConfiguration struct {
	connectionUri  string
	databaseName   string
	collectionName string
}

func NewDatabaseConfiguration(connectionUri, databaseName, collectionName string) databaseConfiguration {
	return databaseConfiguration{
		connectionUri:  connectionUri,
		databaseName:   databaseName,
		collectionName: collectionName,
	}
}
