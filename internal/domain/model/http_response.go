package model

// HTTPResponse represents the response data for HTTP endpoints
type HTTPResponse struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

// NewStatusResponse creates a new HTTPResponse with only status code
func NewStatusResponse(statusCode int) *HTTPResponse {
	return &HTTPResponse{
		StatusCode: statusCode,
	}
}

// NewTextResponse creates a new HTTPResponse with text content
func NewTextResponse(statusCode int, body []byte) *HTTPResponse {
	return &HTTPResponse{
		StatusCode:  statusCode,
		Body:        body,
		ContentType: "text/plain",
	}
}
