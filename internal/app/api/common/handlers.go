package common

import "github.com/gin-gonic/gin"

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
}
