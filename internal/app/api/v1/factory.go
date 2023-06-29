package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/app/api/v1/handlers"
	"git.lothric.net/examples/go/gogin/internal/app/logic"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
	"git.lothric.net/examples/go/gogin/internal/pkg/metrics"
)

// PathHandler defines an API Handler that could attach
// it's underlying handlers to the parent API group.
type PathHandler interface {

	// AttachTo attaches all the underlying path handlers
	// to the parent API group.
	AttachTo(g *gin.RouterGroup) error
}

// HandlerFactory defines a factory that creates
// API Handlers of a particular kind.
type HandlersFactory interface {

	// CreateGistsHandler creates a API Handler
	// that handles all API under '/api/gists'
	CreateGistsHandler() (PathHandler, error)
}

// handlersFactory is a factory that creates API Handlers
// and implements the 'HandlersFactory' interface.
type handlersFactory struct {
	log logger.Log
}

// NewHandlersFactory creates a new API handler factory that
// creates a specific API Handlers.
func NewHandlersFactory(
	ctx context.Context,
	log logger.Log,
) (HandlersFactory, error) {

	return &handlersFactory{
		log: log.WithField(logger.FieldPackage, "v1"),
	}, nil
}

// CreateGistsHandler creates a API Handler
// that handles all API under '/api/gists'
func (f *handlersFactory) CreateGistsHandler() (PathHandler, error) {
	log := f.log.WithField(logger.FieldFunction, "CreateGistsHandler")
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

	handler, err := handlers.NewGistsHandler(f.log, bl)
	if err != nil {
		// TODO: add logger entry
		return nil, err
	}

	return handler, nil
}
