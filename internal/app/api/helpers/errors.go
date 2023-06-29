package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/app/api/v1/models"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

// AbortWithError is a wrapper that simplifies error cases,
// when we return exactly one single error without any detail.
func AbortWithError(
	c *gin.Context,
	log logger.Log,
	statusCode int,
	errorCode string,
	errorMsg string,
) {

	// Report the error to the log
	log.Error(errors.New(errorCode), errorMsg)

	// Abort the gin context and return the single error
	c.AbortWithStatusJSON(
		statusCode,
		[]models.Error{
			{
				Code:    errorCode,
				Message: errorMsg,
			},
		},
	)

}
