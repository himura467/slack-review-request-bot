package model

// HTTPRequest represents the request data for HTTP endpoints
type HTTPRequest struct {
	Body    []byte
	Headers map[string][]string
}

func NewHTTPRequest(body []byte, headers map[string][]string) *HTTPRequest {
	return &HTTPRequest{
		Body:    body,
		Headers: headers,
	}
}
