package ginmw

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WrapperH wraps a typed handler function into a gin.HandlerFunc.
// The handler receives a parsed request and returns a response and error.
// The response is wrapped in a unified APIResponse{Code, Message, Data} format.
//
// On success (err == nil): {code: 0, message: "success", data: <response>}
// On error implementing Coder: {code: coder.Code(), message: coder.Message()}
// On plain error: {code: -1, message: err.Error()}
//
// Request binding: GET uses ShouldBindQuery, other methods use ShouldBindJSON.
func WrapperH[Request, Response any](fn func(*gin.Context, *Request) (*Response, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := getLogger()

		var req Request
		var err error
		if c.Request.Method == http.MethodGet {
			err = c.ShouldBindQuery(&req)
		} else {
			err = c.ShouldBindJSON(&req)
		}

		pk := getPacker()

		if err != nil {
			log.Errorf("parse param err: %v", err)
			c.JSON(http.StatusOK, pk.PackError(err))
			return
		}

		reqBytes, _ := json.Marshal(req)
		log.Infof("[REQ] URI:(%s) BODY:(%s)", c.Request.URL.Path, string(reqBytes))

		rsp, err := fn(c, &req)
		if err != nil {
			log.Errorf("[RSP] URI:(%s) err: %v", c.Request.URL.Path, err)
			c.JSON(http.StatusOK, pk.PackError(err))
			return
		}

		resp := pk.PackSuccess(rsp)
		rspBytes, _ := json.Marshal(resp)
		log.Infof("[RSP] URI:(%s) BODY:(%s)", c.Request.URL.Path, string(rspBytes))
		c.JSON(http.StatusOK, resp)
	}
}
