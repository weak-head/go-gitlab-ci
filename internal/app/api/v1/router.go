package v1

import (
	"errors"

	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

const (
	// gistsRoute is the parent route for gists
	gistsRoute = "gists"
)

var (
	// ErrNoLoggerProvided happens when Logger is not provided.
	ErrNoLoggerProvided = errors.New("no logger provided")

	// ErrNoGistsHandlerProvided happens when Gists Handler is not provided.
	ErrNoGistsHandlerProvided = errors.New("no gists handler provided")
)

// PathHandler defines an API Handler that could attach
// it's underlying handlers to the parent API group.
type PathHandler interface {

	// AttachTo attaches all the underlying path handlers
	// to the parent API group.
	AttachTo(g *gin.RouterGroup) error
}

// v1Router is a v1 API root-level path handler, that constructs all underlying
// API groups that constitute v1 API groups.
type v1Router struct {
	log          logger.Log
	gistsHandler PathHandler
}

// NewV1PathHandler creates a new API v1 root level
// path handler that abstract all underlying v1 API routes.
func NewV1Router(
	log logger.Log,
	gistsHandler PathHandler,
) (PathHandler, error) {

	if log == nil {
		return nil, ErrNoLoggerProvided
	}

	if gistsHandler == nil {
		return nil, ErrNoGistsHandlerProvided
	}

	return &v1Router{
		log:          log.WithField(logger.FieldPackage, "v1"),
		gistsHandler: gistsHandler,
	}, nil
}

// AttachTo attaches all underlying v1 API children to the root level router group.
func (p *v1Router) AttachTo(g *gin.RouterGroup) error {
	log := p.log.WithField(logger.FieldFunction, "AttachTo")
	log.Info("Handling AttachTo")

	// ------------
	// Attach gists handler to the parent API group.
	gistsGroup := g.Group(gistsRoute)
	p.gistsHandler.AttachTo(gistsGroup)

	// ------------
	// Note: Attach more handlers here
	// ------------

	return nil
}
