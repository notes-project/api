package server

import (
	"sync"

	"go.uber.org/zap"
)

type ServerFactory interface {
	NewServer(serverConfig serverConfiguration) Server
}

type serverFactory struct{}

func NewServerFactory() ServerFactory {
	return serverFactory{}
}

var (
	once sync.Once

	serverInstance Server
)

func (sf serverFactory) NewServer(serverConfig serverConfiguration) Server {

	once.Do(func() {
		serverInstance = server{
			serverConfiguration: serverConfig,
			logger:              zap.L().Named("Server"),
		}
	})

	return serverInstance
}
