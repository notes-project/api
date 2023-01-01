package main

import (
	"fmt"
	"os"

	"github.com/denislavPetkov/notes/pkg/database"
	"github.com/denislavPetkov/notes/pkg/server"
	"github.com/denislavPetkov/notes/pkg/utils"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)

	config, err := utils.GetConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get server configuration, err: %s", err.Error()))
		os.Exit(1)
	}

	database := database.NewDatabaseFactory().NewDatabase(config.DatabaseUri, config.DatabaseName, config.DatabaseCollection)

	err = database.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to the database, err: %s", err.Error()))
		os.Exit(1)
	}

	server := server.NewServerFactory().NewServer(config.ServerPort, database)

	server.Start()

}
