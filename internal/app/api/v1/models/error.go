package models

// Error is a single HTTP API request processing error.
//
//	@Description	Error is a single error that has happened during the HTTP API request processing.
//	@Description
//	@Description	Sometimes, we may want to report more than one error for a request.
//	@Description	In this case, we should return several errors in a list.
type Error struct {

	// Code contains the unique error code.
	//
	// The code field should not match the response code.
	// Instead, it should be an error code unique to our application.
	// Generally, there is no convention for the error code field, expect that it be unique.
	//
	// Usually, this field contains only alphanumerics and connecting characters, such as dashes or underscores.
	// For example, 0001, auth-0001 and incorrect-user-pass are canonical examples of error codes.
	Code string `json:"code" bson:"code" binding:"required" example:"x-id-name-missing"`

	// Message is the presentable message to a user.
	//
	// The message portion of the body is usually considered presentable on user interfaces.
	// Therefore, we should translate this title if we support internationalization.
	// So if a client sends a request with an Accept-Language header corresponding to French,
	// the title value should be translated to French.
	Message string `json:"message" bson:"message" binding:"required" example:"The 'x-id-name' header is missing."`

	// Detail is optional information that targets developers and could help trace and investigate
	// the internal reasons why the error has happened.
	//
	// The detail portion is intended for use by developers of clients and not the end user,
	// so the translation is not necessary.
	Detail string `json:"detail" bson:"detail" binding:"required" example:"internal.pkg.api.middleware.idcheck"`
}
