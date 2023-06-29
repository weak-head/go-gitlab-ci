package extensions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*

This is extension to gin that implements API versioning via HTTP headers.

So we have the following routing functionality via providing HTTP header:
> curl --location 'localhost:8082/user' --header 'Accept-version: v1'
> curl --location 'localhost:8082/user' --header 'Accept-version: v2'


Here is an example of how this extension could be used to provide the support,
for API versioning via "Accept-version" HTTP header

 engine := gin.Default()
 router := api.NewRouter(engine)

 defaultRouter := router.Default()
 defaultRouter.Get("/profile",func(ctx *gin.Context) {

 })

 v1 := router.WithVersion("/v1")
 v1.Get("/user",func(ctx *gin.Context) {
  ctx.String(http.StatusOK, "This is the profile v1 API")
 })

 v2 := router.WithVersion("/v2")
 v2.Get("/user",func(ctx *gin.Context) {
  ctx.String(http.StatusOK, "This is the profile v2 API")
 })

*/

type Router struct {
	router        *gin.Engine
	versionGroups map[string]*gin.RouterGroup
}

type VersionedRouter struct {
	version string
	Router
}

func NewRouter(router *gin.Engine) *Router {
	return &Router{
		router:        router,
		versionGroups: make(map[string]*gin.RouterGroup),
	}
}

func (a *Router) Default() VersionedRouter {
	return VersionedRouter{Router: *a}
}

func (a *Router) WithVersion(version string) VersionedRouter {
	if _, ok := a.versionGroups[version]; ok {
		panic("cannot initialize same version multiple times")
	}

	a.versionGroups[version] = a.router.Group(version)
	return VersionedRouter{Router: *a, version: version}
}

func (vr VersionedRouter) Get(relativePath string, handlers ...gin.HandlerFunc) {
	vr.handle(http.MethodGet, relativePath, handlers...)
}

func (vr VersionedRouter) Post(relativePath string, handlers ...gin.HandlerFunc) {
	vr.handle(http.MethodPost, relativePath, handlers...)
}

// Note: You need to follow the same for other HTTP Methods.
// As an example, we can write a method for Put HTTP Method as below,
//
//  func (vr VersionedRouter) Put(relativePath string, handlers ...gin.HandlerFunc)  {
//   vr.handle(http.MethodPost,relativePath,handlers...)
//  }

func (vr VersionedRouter) handle(method, relativePath string, handlers ...gin.HandlerFunc) {
	if !vr.isRouteExist(method, relativePath) {
		vr.router.Handle(method, relativePath, func(ctx *gin.Context) {
			version := ctx.Request.Header.Get("Accept-version")
			if len(version) == 0 {
				ctx.String(http.StatusBadRequest, "Accept-version header is empty")
			}
			ctx.Request.URL.Path = fmt.Sprintf("/%s%s", version, ctx.Request.URL.Path)
			vr.router.HandleContext(ctx)
		})
	}

	versionedRelativePath := vr.version + relativePath
	if !vr.isRouteExist(method, versionedRelativePath) {
		vr.router.Handle(method, versionedRelativePath, handlers...)
	}
}

func (a VersionedRouter) isRouteExist(method, relativePath string) bool {
	for _, route := range a.router.Routes() {
		if route.Method == method && relativePath == route.Path {
			return true
		}
	}
	return false
}
