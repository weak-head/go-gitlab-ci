package constants

const (
	// ErrUnknownErrorCode uniquely identifies the case when we have encountered unknown error
	ErrUnknownErrorCode = "unknown-error"

	// ErrUnknownErrorMsg is returned when the service has encountered unknown error
	ErrUnknownErrorMsg = "The service has encountered unexpected error that it was not able to handle."

	// ErrFailedToParseRequestJsonCode uniquely identifies the cases when we have
	// failed to parse the request content
	ErrFailedToParseRequestJsonCode = "failed-to-parse-request-json"

	// ErrFailedToParseRequestJsonMsg happens when we have failed to parse JSON request content
	ErrFailedToParseRequestJsonMsg = "Failed to parse JSON request content."
)
