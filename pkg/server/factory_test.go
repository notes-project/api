package server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {

	var (
		serverFactory serverFactory = serverFactory{}
	)

	Describe("NewServerFactory", func() {
		It("should return a new server factory object", func() {
			serverFactory := NewServerFactory()
			Expect(serverFactory).NotTo(BeNil())
		})
	})

	Describe("NewServer", func() {
		It("should return a new server object", func() {
			serverInstance := serverFactory.NewServer(serverConfiguration{})
			Expect(serverInstance).NotTo(BeNil())
		})
	})

})
