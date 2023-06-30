package v1

import (
	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/app/api/common"
	"git.lothric.net/examples/go/gogin/internal/app/api/v1/handlers"
	"git.lothric.net/examples/go/gogin/internal/app/logic"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
	"git.lothric.net/examples/go/gogin/internal/pkg/metrics"
)

// v1PathHandler is a API v1 root path handler
// that constructs and handles all underlying
// v1 API groups.
type v1PathHandler struct {
	log     logger.Log
	factory common.ComponentFactory
}

// NewV1PathHandler creates a new API v1 root level
// path handler that abstract all underlying v1 API routes.
func NewV1PathHandler(
	log logger.Log,
	factory common.ComponentFactory,
) (common.PathHandler, error) {
	return &v1PathHandler{
		log:     log.WithField(logger.FieldPackage, "v1"),
		factory: factory,
	}, nil
}

// AttachTo attaches all underlying v1 API children to the root level router group.
func (p *v1PathHandler) AttachTo(g *gin.RouterGroup) error {
	log := p.log.WithField(logger.FieldFunction, "AttachTo")
	log.Info("Handling AttachTo")

	// ------------
	// Create and attach 'gists' handler to the parent API group.
	gistsGroup := g.Group("gists")
	gistsHandler, err := p.createGistsHandler()
	if err != nil {
		log.Error(err, "Failed to create '/gists' path handler")
		return err
	}
	gistsHandler.AttachTo(gistsGroup)

	// ------------
	// Note: Create and attach more handlers here
	// ------------

	return nil
}

func (p *v1PathHandler) createGistsHandler() (common.PathHandler, error) {
	log := p.log.WithField(logger.FieldFunction, "createGistsHandler")
	log.Info("Handling CreateGistsHandler")

	mr, err := metrics.NewReporter()
	if err != nil {
		// TODO: add logger entry
		return nil, err
	}

	bl, err := logic.NewGistsLogic(log, mr)
	if err != nil {
		// TODO: add logger entry
		return nil, err
	}

	handler, err := handlers.NewGistsHandler(log, bl)
	if err != nil {
		// TODO: add logger entry
		return nil, err
	}

	return handler, nil
}
