package ginmw

import "sync"

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

// ResponsePacker defines how WrapperH packs success and error responses.
// Implement this interface to customize the response format.
type ResponsePacker interface {
	PackSuccess(data any) any
	PackError(err error) any
}

var (
	packer   ResponsePacker = defaultPacker{}
	packerMu sync.RWMutex
)

// SetResponsePacker sets the global response packer used by WrapperH.
// If p is nil, the default packer (APIResponse format) is restored.
func SetResponsePacker(p ResponsePacker) {
	packerMu.Lock()
	defer packerMu.Unlock()
	if p == nil {
		packer = defaultPacker{}
		return
	}
	packer = p
}

func getPacker() ResponsePacker {
	packerMu.RLock()
	defer packerMu.RUnlock()
	return packer
}

// defaultPacker is the built-in packer that produces APIResponse.
type defaultPacker struct{}

func (defaultPacker) PackSuccess(data any) any {
	return &APIResponse{Code: 0, Message: "success", Data: data}
}

func (defaultPacker) PackError(err error) any {
	return buildErrorResponse(err)
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
