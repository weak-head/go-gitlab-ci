package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/app/api/constants"
	"git.lothric.net/examples/go/gogin/internal/app/api/helpers"
	"git.lothric.net/examples/go/gogin/internal/app/api/v1/models"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

const (

	// pkg is this package for logger purposes
	pkg = "api.v1.handlers"

	// QueryLanguage is a query key that is used to specify the programming language.
	QueryLanguage = "lang"
)

// GistsLogic is a business logic layer for gists related functionality.
type GistsLogic interface {

	// GetGists returns a filtered list of gists
	GetGists(ctx context.Context, language string) error
}

// MetricsReporter is metrics reporting handler for gists APIs.
type ApiMetricsReporter interface {

	// ApiRequestProcessed
	ApiRequestProcessed(operation string, milliseconds float64)

	// ApiRequestFailed
	ApiRequestFailed(operation string, failure string)
}

// gistsHandler handles all APIs calls for the 'gists' resource.
type gistsHandler struct {
	log     logger.Log
	logic   GistsLogic
	metrics ApiMetricsReporter
}

// NewGistsHandler creates a new instance of the API handler
// that handles all requests to 'gists' resource.
func NewGistsHandler(
	log logger.Log,
	logic GistsLogic,
	metrics ApiMetricsReporter,
) (*gistsHandler, error) {

	gh := &gistsHandler{
		log:     log.WithField(logger.FieldPackage, pkg),
		logic:   logic,
		metrics: metrics,
	}

	return gh, nil
}

// AttachTo attaches the templatesHandler to the
// provided parent router group.
func (gh *gistsHandler) AttachTo(g *gin.RouterGroup) error {

	// GET /api/gists?lang=haskell
	g.GET("", gh.getGists)

	// POST /api/gists
	g.POST("", gh.postGist)

	// GET /api/gists/{id}
	g.GET(":id", gh.getGist)

	// PUT /api/gists/{id}
	g.PUT(":id", gh.putGist)

	// DELETE /api/gists/{id}
	g.DELETE(":id", gh.deleteGist)

	return nil
}

// getGists godoc
//
//	@Summary		Get the list of Source Code Gists, filtered by programming language.
//	@Description	This method returns the list of Gists, that are created using a particular programming language.
//	@Description	This is filtered subset of all available Gists.
//	@Tags			Gists
//	@Param			lang	query	string	false	"Programming language"
//	@Produce		json
//	@Success		200	{array}	models.GistInfo	"The list of Gists has been successfully retrieved."
//	@Failure		500	{array}	models.Error	"The service has encountered unexpected error that it was not able to handle."
//	@Router			/gists [get]
func (gh *gistsHandler) getGists(c *gin.Context) {
	log, ctx, _, err := helpers.ParseContext(gh.log, c, "getGists")
	if err != nil {
		helpers.AbortWithError(c, log,
			http.StatusInternalServerError,
			constants.ErrUnknownErrorCode,
			constants.ErrUnknownErrorMsg)
		return
	}
	log.Info("Handling getGists")

	// Extract argument
	lang := c.Query(QueryLanguage)

	// Note: just an example of the call
	_ = gh.logic.GetGists(ctx, lang)

	gistsInfo := []models.GistInfo{
		{
			Id:          "d17043a0-216c-4c56-9127-b0bf5e3a4c16",
			Language:    lang,
			Name:        "Convert PNG to JPG",
			Description: "Convert PNG image to JPG image using pillow library",
		},
		{
			Id:          "9faa0e0f-e189-4cff-ae6f-7c0c52b3c876",
			Language:    lang,
			Name:        "Generate unique ID",
			Description: "Generate a unique UUIDv4",
		},
	}

	c.AbortWithStatusJSON(http.StatusOK, gistsInfo)
}

// postGist godoc
//
//	@Summary		Create a new Gist.
//	@Description	This method is called to create and store a new Gist
//	@Tags			Gists
//	@Param			gist	body	models.Gist	true	"Gist definition"
//	@Produce		json
//	@Success		201	{object}	models.GistInfo	"Gist has been created."
//	@Failure		400	{array}		models.Error	"Failed to parse JSON request content."
//	@Failure		500	{array}		models.Error	"The service has encountered unexpected error that it was not able to handle."
//	@Router			/gists [post]
func (gh *gistsHandler) postGist(c *gin.Context) {
	log, _, _, err := helpers.ParseContext(gh.log, c, "postGist")
	if err != nil {
		helpers.AbortWithError(c, log,
			http.StatusInternalServerError,
			constants.ErrUnknownErrorCode,
			constants.ErrUnknownErrorMsg)
		return
	}
	log.Info("Handling postGist")

	// Extract argument
	var gist models.Gist
	if err := c.BindJSON(&gist); err != nil {
		helpers.AbortWithError(c, log,
			http.StatusBadRequest,
			constants.ErrFailedToParseRequestJsonCode,
			constants.ErrFailedToParseRequestJsonMsg)
		return
	}

	gistInfo := models.GistInfo{
		Id:          "d17043a0-216c-4c56-9127-b0bf5e3a4c16",
		Language:    "haskell",
		Name:        "Convert PNG to JPG",
		Description: "Convert PNG image to JPG image using pillow library",
	}

	// Return result
	c.AbortWithStatusJSON(http.StatusCreated, gistInfo)
}

// getGist godoc
//
//	@Summary	Get the detailed information about the Gist.
//	@Tags		Gists
//	@Param		id	path	string	true	"Gist id"
//	@Produce	json
//	@Success	200	{object}	models.GistDetails	"The Gist definition has been successfully retrieved."
//	@Failure	404	{array}		models.Error		"The specified Gist does not exist."
//	@Failure	500	{array}		models.Error		"The service has encountered unexpected error that it was not able to handle."
//	@Router		/gists/{id} [get]
func (gh *gistsHandler) getGist(c *gin.Context) {
	log, _, _, err := helpers.ParseContext(gh.log, c, "getGist")
	if err != nil {
		helpers.AbortWithError(c, log,
			http.StatusInternalServerError,
			constants.ErrUnknownErrorCode,
			constants.ErrUnknownErrorMsg)
		return
	}
	log.Info("Handling getGist")

	// Extract argument
	id := c.Param("id")

	gistDetails := models.GistDetails{
		Id:           id,
		Language:     "haskell",
		Name:         "Convert PNG to JPG",
		Code:         "for (let i = 0; i < 5; i++) {...}",
		Description:  "Convert PNG image to JPG image using pillow library",
		CreatedAt:    "2023-06-07T18:27:25-04:00",
		LastAccessed: "2023-06-24T08:13:59-04:00",
		LastUpdated:  "2023-06-11T10:44:17-04:00",
	}

	// Return result
	c.AbortWithStatusJSON(http.StatusOK, gistDetails)
}

// putGist godoc
//
//	@Summary		Create or replace the Gist.
//	@Description	This method is called to update and store an existing Gist definition.
//	@Description
//	@Tags		Gists
//	@Param		id			path	string		true	"Gist id"
//	@Param		template	body	models.Gist	true	"Gist definition"
//	@Produce	json
//	@Success	201	{object}	models.GistInfo	"Gist has been updated."
//	@Failure	400	{array}		models.Error	"Failed to parse JSON request content."
//	@Failure	500	{array}		models.Error	"The service has encountered unexpected error that it was not able to handle."
//	@Router		/gists/{id} [put]
func (gh *gistsHandler) putGist(c *gin.Context) {
	log, _, _, err := helpers.ParseContext(gh.log, c, "putGist")
	if err != nil {
		helpers.AbortWithError(c, log,
			http.StatusInternalServerError,
			constants.ErrUnknownErrorCode,
			constants.ErrUnknownErrorMsg)
		return
	}
	log.Info("Handling putGist")

	// Extract argument
	id := c.Param("id")
	var gist models.Gist
	if err := c.BindJSON(&gist); err != nil {
		helpers.AbortWithError(c, log,
			http.StatusBadRequest,
			constants.ErrFailedToParseRequestJsonCode,
			constants.ErrFailedToParseRequestJsonMsg)
		return
	}

	gistInfo := models.GistInfo{
		Id:          id,
		Language:    gist.Language,
		Name:        gist.Name,
		Description: gist.Description,
	}

	// Return result
	c.AbortWithStatusJSON(http.StatusCreated, gistInfo)
}

// deleteGist godoc
//
//	@Summary		Delete previously created Gist.
//	@Description	This method is called to delete an existing Gist definition.
//	@Tags			Gists
//	@Param			id	path	string	true	"Gist id"
//	@Produce		json
//	@Success		204	{object}	models.GistInfo	"Gist has been deleted."
//	@Failure		404	{array}		models.Error	"The specified Gist does not exist."
//	@Failure		500	{array}		models.Error	"The service has encountered unexpected error that it was not able to handle."
//	@Router			/gists/{id} [delete]
func (gh *gistsHandler) deleteGist(c *gin.Context) {
	log, _, _, err := helpers.ParseContext(gh.log, c, "deleteGist")
	if err != nil {
		helpers.AbortWithError(c, log,
			http.StatusInternalServerError,
			constants.ErrUnknownErrorCode,
			constants.ErrUnknownErrorMsg)
		return
	}
	log.Info("Handling deleteGist")

	// Extract argument
	id := c.Param("id")

	gistInfo := models.GistInfo{
		Id:          id,
		Language:    "haskell",
		Name:        "Convert PNG to JPG",
		Description: "Convert PNG image to JPG image using pillow library",
	}

	// Return result
	c.AbortWithStatusJSON(http.StatusNoContent, gistInfo)
}
