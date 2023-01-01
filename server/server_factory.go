package server

import (
	"sync"

	"github.com/denislavPetkov/notes/database"
	"go.uber.org/zap"
)

type ServerFactory interface {
	NewServer(port string, db database.Database) Server
}

type serverFactory struct{}

func NewServerFactory() ServerFactory {
	return serverFactory{}
}

var (
	once sync.Once

	serverInstance Server
)

func (sf serverFactory) NewServer(port string, db database.Database) Server {

	once.Do(func() {
		serverInstance = server{
			port:   port,
			db:     db,
			logger: zap.L().Named("Server"),
		}
	})

	return serverInstance
}
