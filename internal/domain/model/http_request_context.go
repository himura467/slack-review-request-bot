package model

// HTTPRequestContext represents the context of an HTTP request
type HTTPRequestContext struct {
	Body    []byte
	Headers map[string][]string
}

func NewHTTPRequestContext(body []byte, headers map[string][]string) *HTTPRequestContext {
	return &HTTPRequestContext{
		Body:    body,
		Headers: headers,
	}
}
