package api

import (
	"context"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"git.lothric.net/examples/go/gogin/internal/app/api/common"
	"git.lothric.net/examples/go/gogin/internal/app/api/middleware"
	v1 "git.lothric.net/examples/go/gogin/internal/app/api/v1"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

// BuildApiEngine creates the API Engine that handles all
// HTTP REST calls to the Gogin APIs.
//
// The API Engine handles the routing, authentication
// and cookie/header validation logic.
func BuildApiEngine(
	ctx context.Context,
	log logger.Log,
	factory common.ComponentFactory,
) (*gin.Engine, error) {

	// We use the vanilla 'gin' router that doesn't support routing of APIs
	// based on "Accept-version" HTTP header.
	//
	// In the far future we could migrate to "v2 API", and we could consider
	// using HTTP Header based routing.
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

	// All our APIs are under 'api' group and require auth
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

	// v1 PathHandler
	v1routes, err := v1.NewV1PathHandler(log, factory)
	if err != nil {
		log.Error(err, "Failed to create v1 PathHandler")
		return nil, err
	}

	// v1 API group is directly attached to the "/api" root
	// so all routes are under "/api".
	// for example for 'templates' the route is:
	//   -> /api/templates
	v1routes.AttachTo(apiGroup)

	// Handle requests to swagger under '/swagger/*'
	// http://localhost:8080/swagger/index.html
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return engine, nil
}
