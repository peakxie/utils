package ginmw

// APIResponse is the unified response format returned by WrapperH handlers.
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Coder is an interface that errors can implement to provide
// a business error code and message.
type Coder interface {
	Code() int
	Message() string
}

func buildErrorResponse(err error) *APIResponse {
	if c, ok := err.(Coder); ok {
		return &APIResponse{
			Code:    c.Code(),
			Message: c.Message(),
		}
	}
	return &APIResponse{
		Code:    -1,
		Message: err.Error(),
	}
}
