package helpers

import (
	"context"
	"errors"

	"git.lothric.net/examples/go/gogin/internal/app/api/constants"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

var (

	// ErrFailedToParseGinContext happens when gin context parsing has failed
	ErrFailedToParseGinContext = errors.New("failed to parse gin context")
)

// ContextModel represents the most important fields
// that are extracted from HTTP headers and JWT Auth token
// during the HTTP Request processing by our defined middleware layer.
//
// This model includes ONLY fields that are extracted from either:
//   - HTTP Request Headers
//   - HTTP Request Cookies
//   - HTTP Request JWT Auth Token
//
// This model doesn't include and SHOULD NEVER INCLUDE:
//   - HTTP Query string parameters
//   - HTTP Body values or arguments
//
// Because those are part of API definitions and should be
// handled by the corresponding HTTP API endpoint handler.
type ContextModel struct {

	// CorrelationId is a unique id that helps to establish
	// the unique correlation of inter-service API calls
	// as well as internal calls between components of the system
	CorrelationId string

	// TODO: add more fields here
}

// ExtractContextModel extract ContextModel from the
// provided gin.Context based on the fields
// extracted by the middleware layer
func ExtractContextModel(c *gin.Context) (ContextModel, error) {

	// Extract fields from the gin context
	corrId := c.GetString(constants.HeaderCorrelationId)

	// Create context model
	m := ContextModel{
		CorrelationId: corrId,
	}

	// Note: return ErrFailedToParseGinContext defined above
	// if parsing or validation has failed

	return m, nil
}

// ParseContext parses gin.Context and extract from it:
//   - [logger.Log]
//   - [context.Context]
//   - [ContextModel]
//
// ParseContext should never fail if there are not bugs in
// the middleware layer.
//
// If ParseContext fails it means we have encountered:
// `http.StatusInternalServerError`, because by this point
// our middleware layer has handled the:
//   - cookies validation
//   - HTTP headers validation
//   - authorization verification
//
// and we have not started yet parsing and validation of:
//   - query parameters
//   - json body
//
// so it means that we have a bug somewhere in the
// logic of middleware layer and something wrong with it.
// So if ParseContext fails, it is always HTTP status code 500
// and the signal to find and fix the bug.
//
// Note: the log is always correctly initialized
// Note: and returned from this method, even if an error is returned
func ParseContext(
	parentLog logger.Log,
	c *gin.Context,
	function string,
) (logger.Log, context.Context, ContextModel, error) {

	// Create logger with additional fields, that are
	// extracted from the gin.Context
	log, ctx := logger.FromGin(parentLog, c, function)

	// Extract ContextModel from the gin.Context
	cm, err := ExtractContextModel(c)

	// This is our own error that we have created
	// because we were not able to extract some
	// information from the context that we have been
	// expecting to be there.
	if err != nil {
		log.Error(err, "Failed to parse the gin context")
		return log, nil, ContextModel{}, err
	}

	return log, ctx, cm, nil
}
