package ginmw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peakxie/utils/internal/log"
	"github.com/peakxie/utils/loghelper"
)

// WrapperH ...
func WrapperH[Request, Response any](fn func(*gin.Context, *Request) (*Response, error)) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req Request
		var err error
		if c.Request.Method == "GET" {
			err = c.ShouldBindQuery(&req)
		} else {
			err = c.ShouldBindJSON(&req)
		}

		if err != nil {
			log.Errorf("parse param err: %v", err)
			return
		}

		log.Infof("[REQ] URI:(%s) BODY:(%s)", c.Request.URL.Path, loghelper.ToPrintString(req))

		rsp, err := fn(c, &req)
		if err != nil {
			log.Errorf("rsp err: %v", err)
			return
		}

		c.JSON(http.StatusOK, rsp)
	}
}
