package utils

import (
	"fmt"
	"os"
)

const (
	DATABASE_URI        = "DATABASE_URI"
	DATABASE_NAME       = "DATABASE_NAME"
	DATABASE_COLLECTION = "DATABASE_COLLECTION"

	SERVER_PORT = "SERVER_PORT"
)

const (
	envVarIsEmptyErrMsg = "env var %s is empty"
)

type ServerConfig struct {
	DatabaseUri        string
	DatabaseName       string
	DatabaseCollection string

	ServerPort string
}

func GetServerConfig() (ServerConfig, error) {
	dbUri, exist := os.LookupEnv(DATABASE_URI)
	if !exist {
		return ServerConfig{}, fmt.Errorf(envVarIsEmptyErrMsg, DATABASE_URI)
	}

	dbName, exist := os.LookupEnv(DATABASE_NAME)
	if !exist {
		return ServerConfig{}, fmt.Errorf(envVarIsEmptyErrMsg, DATABASE_NAME)
	}

	dbCollection, exist := os.LookupEnv(DATABASE_COLLECTION)
	if !exist {
		return ServerConfig{}, fmt.Errorf(envVarIsEmptyErrMsg, DATABASE_COLLECTION)
	}

	serverPort, exist := os.LookupEnv(SERVER_PORT)
	if !exist {
		return ServerConfig{}, fmt.Errorf(envVarIsEmptyErrMsg, serverPort)
	}

	return ServerConfig{
		DatabaseUri:        dbUri,
		DatabaseName:       dbName,
		DatabaseCollection: dbCollection,
		ServerPort:         serverPort,
	}, nil
}
