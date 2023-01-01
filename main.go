package main

import (
	"fmt"
	"os"

	"github.com/denislavPetkov/notes/database"
	"github.com/denislavPetkov/notes/server"
	"github.com/denislavPetkov/notes/utils"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)

	serverConfig, err := utils.GetServerConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get server configuration, err: %s", err.Error()))
		os.Exit(1)
	}

	database := database.NewDatabaseFactory().NewDatabase(serverConfig.DatabaseUri, serverConfig.DatabaseName, serverConfig.DatabaseCollection)

	err = database.Start()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to the database, err: %s", err.Error()))
		os.Exit(1)
	}

	server := server.NewServerFactory().NewServer(serverConfig.ServerPort, database)

	server.Start()

}
