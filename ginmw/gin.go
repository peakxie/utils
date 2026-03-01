package ginmw

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// LogFunc is a context-aware printf-style logging function.
// The *gin.Context parameter allows callers to extract per-request
// information (e.g. trace ID, request ID) when logging.
type LogFunc func(c *gin.Context, format string, args ...any)

var (
	logFn      LogFunc
	errorLogFn LogFunc
	logFnMu    sync.RWMutex
)

// SetLogFunc sets the global log function used by WrapperH.
// If fn is nil, logging is disabled.
func SetLogFunc(fn LogFunc) {
	logFnMu.Lock()
	defer logFnMu.Unlock()
	logFn = fn
}

// SetErrorLogFunc sets the global error-level log function used by WrapperH.
// When set, error logs (param parse failures, handler errors) use this function
// instead of the regular LogFunc. If fn is nil, error logs fall back to LogFunc.
func SetErrorLogFunc(fn LogFunc) {
	logFnMu.Lock()
	defer logFnMu.Unlock()
	errorLogFn = fn
}

func getLogFunc() LogFunc {
	logFnMu.RLock()
	defer logFnMu.RUnlock()
	return logFn
}

func getErrorLogFunc() LogFunc {
	logFnMu.RLock()
	defer logFnMu.RUnlock()
	if errorLogFn != nil {
		return errorLogFn
	}
	return logFn
}

// WrapperH wraps a typed handler function into a gin.HandlerFunc.
// The handler receives a parsed request and returns a response and error.
// The response is wrapped via the global ResponsePacker (default: APIResponse{Code, Message, Data}).
//
// Request/response logging uses the global LogFunc set via SetLogFunc.
// If no LogFunc is set, no logging is performed.
//
// On success (err == nil): {code: 0, message: "success", data: <response>}
// On error implementing Coder: {code: coder.Code(), message: coder.Message()}
// On plain error: {code: -1, message: err.Error()}
//
// Request binding: GET uses ShouldBindQuery, other methods use ShouldBindJSON.
func WrapperH[Request, Response any](fn func(*gin.Context, *Request) (*Response, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := getLogFunc()
		errLog := getErrorLogFunc()

		var req Request
		var err error
		if c.Request.Method == http.MethodGet {
			err = c.ShouldBindQuery(&req)
		} else {
			err = c.ShouldBindJSON(&req)
		}

		pk := getPacker()

		if err != nil {
			if errLog != nil {
				errLog(c, "parse param err: %v", err)
			}
			c.JSON(http.StatusOK, pk.PackError(err))
			return
		}

		if log != nil {
			reqBytes, _ := json.Marshal(req)
			log(c, "[REQ] URI:(%s) BODY:(%s)", c.Request.URL.Path, string(reqBytes))
		}

		rsp, err := fn(c, &req)
		if err != nil {
			if errLog != nil {
				errLog(c, "[RSP] URI:(%s) err: %v", c.Request.URL.Path, err)
			}
			c.JSON(http.StatusOK, pk.PackError(err))
			return
		}

		resp := pk.PackSuccess(rsp)
		if log != nil {
			rspBytes, _ := json.Marshal(resp)
			log(c, "[RSP] URI:(%s) BODY:(%s)", c.Request.URL.Path, string(rspBytes))
		}
		c.JSON(http.StatusOK, resp)
	}
}
