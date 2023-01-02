package utils

import (
	"fmt"
	"os"
)

const (
	DATABASE_URI        = "DATABASE_URI"
	DATABASE_NAME       = "DATABASE_NAME"
	DATABASE_COLLECTION = "DATABASE_COLLECTION"

	SERVER_PORT     = "SERVER_PORT"
	SERVER_TLS_PORT = "SERVER_TLS_PORT"
)

const (
	envVarIsEmptyErrMsg = "env var %s is empty"
)

type Config struct {
	DatabaseUri        string
	DatabaseName       string
	DatabaseCollection string

	ServerPort    string
	ServerTlsPort string
}

func GetEnvConfig() (Config, error) {
	dbUri, exist := os.LookupEnv(DATABASE_URI)
	if !exist {
		return Config{}, fmt.Errorf(envVarIsEmptyErrMsg, DATABASE_URI)
	}

	dbName, exist := os.LookupEnv(DATABASE_NAME)
	if !exist {
		return Config{}, fmt.Errorf(envVarIsEmptyErrMsg, DATABASE_NAME)
	}

	dbCollection, exist := os.LookupEnv(DATABASE_COLLECTION)
	if !exist {
		return Config{}, fmt.Errorf(envVarIsEmptyErrMsg, DATABASE_COLLECTION)
	}

	serverPort, exist := os.LookupEnv(SERVER_PORT)
	if !exist {
		return Config{}, fmt.Errorf(envVarIsEmptyErrMsg, serverPort)
	}

	serverTlsPort := os.Getenv(SERVER_TLS_PORT)

	return Config{
		DatabaseUri:        dbUri,
		DatabaseName:       dbName,
		DatabaseCollection: dbCollection,
		ServerPort:         serverPort,
		ServerTlsPort:      serverTlsPort,
	}, nil
}
