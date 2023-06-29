package v1

import (
	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

// v1PathHandler is a API v1 root path handler
// that constructs and handles all underlying
// v1 API groups.
type v1PathHandler struct {
	log     logger.Log
	factory HandlersFactory
}

// NewV1PathHandler creates a new API v1 root level
// path handler that abstract all underlying v1 API routes.
func NewV1PathHandler(
	log logger.Log,
	factory HandlersFactory,
) (PathHandler, error) {
	return &v1PathHandler{
		log:     log.WithField(logger.FieldPackage, "v1"),
		factory: factory,
	}, nil
}

// AttachTo attaches all underlying v1 API children to the
// root level router group.
func (rb *v1PathHandler) AttachTo(g *gin.RouterGroup) error {
	log := rb.log.WithField(logger.FieldFunction, "AttachTo")
	log.Info("Handling AttachTo")

	// Create and attach 'gists' handler
	// to the parent API group.
	gg := g.Group("gists")
	// tg.Use(middleware.AuthenticationRequired(log))
	gh, err := rb.factory.CreateGistsHandler()
	if err != nil {
		// TODO: add log entry
		return err
	}
	gh.AttachTo(gg)

	// Note: Create and attach more administration handlers here

	return nil
}
