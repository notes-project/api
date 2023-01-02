package server

import "github.com/denislavPetkov/notes/pkg/database"

type serverConfiguration struct {
	port string
	db   database.Database

	tlsPort         string
	tlsCertLocation string
	tlsKeyLocation  string
}

func NewServerConfiguration(port, tlsPort, tlsCertLocation, tlsKeyLocation string, db database.Database) serverConfiguration {
	return serverConfiguration{
		port:            port,
		db:              db,
		tlsPort:         tlsPort,
		tlsCertLocation: tlsCertLocation,
		tlsKeyLocation:  tlsKeyLocation,
	}
}
