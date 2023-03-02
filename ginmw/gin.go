package ginmw

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	"github.com/peakxie/utils/loghelper"
)

// fun(c *gin.Context, req interface{}) (interface{}, error)
// reqfun() (interface{})
// rspfun(c *gin.Context, err error) (interface{})
func Wrapper(fun interface{}, reqfun interface{}, rspfun interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

		var rsp interface{}

		defer func() {
			if err := recover(); err != nil {
				rspV := reflect.ValueOf(rspfun).Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(err)})
				if len(rspV) == 1 {
					rsp = rspV[0].Interface()
				}
			}
			fmt.Printf("[RSP] URI:(%s) BODY:(%s)", c.Request.URL.Path, loghelper.ToPrintString(rsp))

			if _, ok := rsp.(proto.Message); ok {
				c.ProtoBuf(http.StatusOK, rsp)
			} else {
				//c.Header("Content-Type", "application/json")
				c.JSON(http.StatusOK, rsp)
			}
		}()
		reqV := reflect.ValueOf(reqfun).Call([]reflect.Value{})
		if len(reqV) != 1 || reqV[0].IsNil() {
			fmt.Errorf("get req error.")
			panic(errors.New("req struct err"))
		}

		req := reqV[0].Interface()
		var err error
		if _, ok := req.(proto.Message); ok {
			err = c.ShouldBindWith(reqV[0].Interface(), binding.ProtoBuf)
		} else {
			err = c.ShouldBindJSON(reqV[0].Interface())
		}
		if err != nil {
			fmt.Errorf("parse param err:%v", err)
			panic(err)
		}

		fmt.Printf("[REQ] URI:(%s) BODY:(%s)", c.Request.URL.Path, loghelper.ToPrintString(reqV[0].Interface()))

		rspVV := reflect.ValueOf(fun).Call([]reflect.Value{reflect.ValueOf(c), reqV[0]})
		if err, ok := rspVV[1].Interface().(error); ok {
			fmt.Errorf("process err: %s", err.Error())
			panic(err)
		}
		rsp = rspVV[0].Interface()
	}
}
