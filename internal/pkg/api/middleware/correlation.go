package middleware

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/pkg/api/constants"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

const (

	// corrIdLetters defines the characters that are used to generate
	// the unique correlation ID in case if the "X-Request-Id" HTTP
	// header is missing and we need to generate the ID ourselves.
	corrIdLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// rndCorrelationId generates a new unique correlation id
func rndCorrelationId() string {
	n := 15 // id len
	b := make([]byte, n)
	for i := range b {
		b[i] = corrIdLetters[rand.Int63()%int64(len(corrIdLetters))]
	}
	return string(b)
}

// EnsureCorrelationId middleware ensures that the Correlation Id
// exist as part of a HTTP request header so, we can track our call chain
// in the log files down to the database level and establish a clear
// connection of calls and events.
//
// Ideally this correlation id should be generated on the
// API Gateway level for each incoming HTTP request from external
// Web, Mobile or Application clients and passed in the "X-Request-Id" header,
// so we can use this correlation id to track the entire
// call chain between several microservices and establish reliable
// tracing and logging for the entire Digital Platform.
//
// In case if correlation id is not provided in the "X-Request-Id"
// header we are generating a new one in this middleware.
//
// More on the 'X-Request-Id' header could be found here:
//   - https://http.dev/x-request-id
func EnsureCorrelationId(log logger.Log) gin.HandlerFunc {

	// Create a closure to capture the adjusted log
	log = log.WithFields(logger.Fields{
		logger.FieldPackage:  "middleware",
		logger.FieldFunction: "EnsureCorrelationId",
	})

	return func(c *gin.Context) {

		log.Info("Verifying correlation header exist")

		// If 'X-Request-Id' header is missing,
		// we generate it and set it
		corrHeader := http.CanonicalHeaderKey(constants.HeaderCorrelationId)
		if c.Request.Header[corrHeader] == nil {

			corrId := rndCorrelationId()
			log.Warnf("Correlation header doesn't exist, generating a new one: [%s]", corrId)

			c.Request.Header.Set(corrHeader, corrId)
		}

		// Set 'x-request-id' to the gin context
		// so we can extract it in the HTTP API handlers
		c.Set(constants.HeaderCorrelationId, c.Request.Header.Get(corrHeader))

		c.Next()
	}
}
