package utils

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	Describe("GetEnvConfig", func() {

		BeforeEach(func() {
			Expect(os.Setenv(DATABASE_URI, "dbUri")).To(Succeed())
			Expect(os.Setenv(DATABASE_NAME, "dbName")).To(Succeed())
			Expect(os.Setenv(DATABASE_COLLECTION, "dbCollection")).To(Succeed())
			Expect(os.Setenv(SERVER_PORT, "dbPort")).To(Succeed())
			Expect(os.Setenv(SERVER_TLS_PORT, "dbTlsPort")).To(Succeed())
		})

		Context("Required env vars", func() {
			It("should return an error when database uri is missing", func() {
				os.Unsetenv(DATABASE_URI)

				_, err := GetEnvConfig()
				Expect(err).To(MatchError(fmt.Sprintf(envVarIsEmptyErrMsg, DATABASE_URI)))
			})

			It("should return an error when database name is missing", func() {
				os.Unsetenv(DATABASE_NAME)

				_, err := GetEnvConfig()
				Expect(err).To(MatchError(fmt.Sprintf(envVarIsEmptyErrMsg, DATABASE_NAME)))
			})

			It("should return an error when database collection is missing", func() {
				os.Unsetenv(DATABASE_COLLECTION)

				_, err := GetEnvConfig()
				Expect(err).To(MatchError(fmt.Sprintf(envVarIsEmptyErrMsg, DATABASE_COLLECTION)))
			})

			It("should return an error when server port is missing", func() {
				os.Unsetenv(SERVER_PORT)

				_, err := GetEnvConfig()
				Expect(err).To(MatchError(fmt.Sprintf(envVarIsEmptyErrMsg, SERVER_PORT)))
			})

			It("should not return an error when server tls port is missing", func() {
				os.Unsetenv(SERVER_TLS_PORT)

				_, err := GetEnvConfig()
				Expect(err).NotTo(HaveOccurred())
			})
		})

	})

})

var _ = AfterSuite(func() {
	Expect(os.Unsetenv(DATABASE_URI)).To(Succeed())
	Expect(os.Unsetenv(DATABASE_NAME)).To(Succeed())
	Expect(os.Unsetenv(DATABASE_COLLECTION)).To(Succeed())
	Expect(os.Unsetenv(SERVER_PORT)).To(Succeed())
	Expect(os.Unsetenv(SERVER_TLS_PORT)).To(Succeed())
})
