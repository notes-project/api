package server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var (
		testPort            = "testPort"
		testTlsPort         = "testTlsPort"
		testTlsCertLocation = "testTlsCertLocation"
		testTlsKeyLocation  = "testTlsKeyLocation"
	)

	Describe("NewServerConfiguration", func() {
		It("should return a new server configuration object", func() {
			serverConfig := NewServerConfiguration(testPort, testTlsPort, testTlsCertLocation, testTlsKeyLocation, nil)

			Expect(serverConfig).To(Equal(
				serverConfiguration{
					port:            testPort,
					tlsPort:         testTlsPort,
					tlsCertLocation: testTlsCertLocation,
					tlsKeyLocation:  testTlsKeyLocation,
				},
			))
		})
	})

})
