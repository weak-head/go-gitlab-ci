package api

import (
	"context"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"git.lothric.net/examples/go/gogin/internal/app/api/middleware"
	v1 "git.lothric.net/examples/go/gogin/internal/app/api/v1"
	"git.lothric.net/examples/go/gogin/internal/app/api/v1/handlers"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

// PathHandler defines an API Handler that could attach
// it's underlying handlers to the parent API group.
type PathHandler interface {

	// AttachTo attaches all the underlying path handlers
	// to the parent API group.
	AttachTo(g *gin.RouterGroup) error
}

// ComponentFactory is responsible for creating components that
// an API handler require to process an incoming API request.
type ComponentFactory interface {

	// CreateMetricsReporter
	CreateMetricsReporter() (handlers.MetricsReporter, error)

	// CreateGistsLogic
	CreateGistsLogic() (handlers.GistsLogic, error)
}

// apiBuilder
type apiBuilder struct {
	log     logger.Log
	factory ComponentFactory
}

// NewApiBuilder
func NewApiBuilder(
	log logger.Log,
	factory ComponentFactory,
) (*apiBuilder, error) {

	return &apiBuilder{
		log:     log,
		factory: factory,
	}, nil
}

// BuildApiEngine creates the API Engine that handles all
// HTTP REST calls to the Gogin APIs.
//
// The API Engine handles the routing, authentication
// and cookie/header validation logic.
func (b *apiBuilder) BuildApi(
	ctx context.Context,
) (*gin.Engine, error) {
	log := b.log.WithField(logger.FieldFunction, "BuildApi")
	log.Info("Building GoGin API")

	// We use the vanilla 'gin' router that doesn't support routing of APIs
	// based on "Accept-version" HTTP header.
	//
	// When we migrate to "v2 API", we could consider using HTTP Header based routing.
	//
	// For a very simple skeleton of a router that can route requests based
	// on the "Accept-version" HTTP header, refer to the 'api/extensions/apiver.go'
	//
	// In case if we go in a direction of using path-based API versioning,
	// and include 'v1' & 'v2' into API path, the vanilla 'gin' could be used as-is
	// and the changes are very quick and simple.
	// But if we go with the HTTP header based API versioning (that is recommended
	// approach for the enterprise-grade software), we would need to replace the 'gin.Default()'
	// router with something like mentioned above that supports header-based routing.
	engine := gin.Default()

	// All our APIs are under 'api' group
	apiGroup := engine.Group("api")

	// Middlewares are executed in the exact order
	// of how they are defined here. The order of
	// middleware execution is important, so keep
	// this in mind if/when a new middleware is added.
	//
	// We isolate Authentication middleware
	// on a particular path level, because some of the
	// routes requires Authentication check
	// and some of then don't require it.
	//
	// These are global root level middlewares
	// that are applied to all path handlers.
	apiGroup.Use(middleware.EnsureCorrelationId(log))

	// v1 API group is directly attached to the "/api" root
	// so all routes are under "/api".
	// for example for 'templates' the route is:
	//   -> /api/templates
	v1router, err := b.buildV1Api()
	if err != nil {
		log.Error(err, "Failed to create v1 API router")
		return nil, err
	}
	v1router.AttachTo(apiGroup)

	// Handle requests to swagger under '/swagger/*'
	// http://localhost:8080/swagger/index.html
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return engine, nil
}

// buildV1Api creates V1 API group
func (b *apiBuilder) buildV1Api() (PathHandler, error) {
	log := b.log.WithField(logger.FieldFunction, "createGistsHandler")
	log.Info("Building v1 API")

	// Shared metrics reporter
	metricsReporter, err := b.factory.CreateMetricsReporter()
	if err != nil {
		log.Error(err, "Failed to create Metrics Reporter")
		return nil, err
	}

	// Gists business logic
	gistsLogic, err := b.factory.CreateGistsLogic()
	if err != nil {
		log.Error(err, "Failed to create Gists Logic")
		return nil, err
	}

	// Gists API handler
	gistsHandler, err := handlers.NewGistsHandler(log, gistsLogic, metricsReporter)
	if err != nil {
		log.Error(err, "Failed to create Gists Handler")
		return nil, err
	}

	// V1 router
	v1router, err := v1.NewV1Router(log, gistsHandler)
	if err != nil {
		log.Error(err, "Failed to create v1 router")
		return nil, err
	}

	return v1router, nil
}
