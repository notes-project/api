package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/denislavPetkov/notes/pkg/database"
	"github.com/denislavPetkov/notes/pkg/server"
	"github.com/denislavPetkov/notes/pkg/utils"
	"go.uber.org/zap"
)

func main() {

	tlsCertLocation := flag.String("tlsCertLocation", "", "Specify location of certificate when TLS is enabled")
	tlsKeyLocation := flag.String("tlsKeyLocation", "", "Specify location of certificate key when TLS is enabled")

	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)

	envConfig, err := utils.GetEnvConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get server configuration, err: %s", err.Error()))
		os.Exit(1)
	}

	dbConfig := database.NewDatabaseConfiguration(envConfig.DatabaseUri, envConfig.DatabaseName, envConfig.DatabaseCollection)

	database := database.NewDatabaseFactory().NewDatabase(dbConfig)

	err = database.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to the database, err: %s", err.Error()))
		os.Exit(1)
	}

	serverConfig := server.NewServerConfiguration(envConfig.ServerPort, envConfig.ServerTlsPort, *tlsCertLocation, *tlsKeyLocation, database)

	server := server.NewServerFactory().NewServer(serverConfig)

	err = server.Start()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start server, err: %s", err.Error()))
		os.Exit(1)
	}

}
