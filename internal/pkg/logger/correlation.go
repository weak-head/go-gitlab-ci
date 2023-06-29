package logger

import (
	"context"

	"github.com/gin-gonic/gin"
)

// context key type for CorrelationID
type contextKeyCorrelationID string

const (

	// HttpCorrelationId is the HTTP header that helps to establish
	// the unique correlation of inter-service API calls
	// as well as internal calls between components of the system
	HeaderCorrelationId = "x-request-id"

	// CorrelationId is the HTTP header that helps to establish
	// the unique correlation of inter-service API calls
	// as well as internal calls between components of the system
	CorrelationId contextKeyCorrelationID = HeaderCorrelationId
)

// FromGin extracts correlationId from Gin context
// and creates a logger and context that has the associated
// correlation id fields set and configured
func FromGin(
	parentLog Log,
	c *gin.Context,
	function string,
) (Log, context.Context) {

	// Extract fields from the gin context
	corrId := c.GetString(HeaderCorrelationId)

	// Logger with the required fields set
	log := parentLog.WithFields(Fields{
		FieldFunction:    function,
		FieldCorrelation: corrId,
	})

	// Create a new context that could be used by other methods,
	// down the call chain with the correlation id set
	ctx := context.WithValue(c.Request.Context(), CorrelationId, corrId)

	return log, ctx
}

// FromContext extracts correlationId from context
// and creates a logger that has the associated
// correlation id field set and configured
func FromContext(
	parentLog Log,
	ctx context.Context,
	function string,
) Log {

	// Extract fields from Context
	corrId := ctx.Value(CorrelationId)

	// Logger with the fields set
	log := parentLog.WithFields(Fields{
		FieldFunction:    function,
		FieldCorrelation: corrId,
	})

	return log
}
