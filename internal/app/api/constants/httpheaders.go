package constants

const (
	// HttpAuthToken is the HTTP Header field that holds JWT authorization token
	//
	// Authorization: Bearer eyJhbGciOiJSUzI1NiIsIn...
	HeaderAuthToken = "Authorization"

	// HttpCorrelationId is the HTTP header that helps to establish
	// the unique correlation of inter-service API calls
	// as well as internal calls between components of the system
	HeaderCorrelationId = "x-request-id"
)
