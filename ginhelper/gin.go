package ginhelper

import (
	"errors"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	"github.com/peakxie/utils/loghelper"
	"github.com/sirupsen/logrus"
)

const (
	REQUESTID string = "_requestid"
	LOGCTX    string = "_logctx"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func() {
		//}()
		id := rand.Uint64()
		ctx.Set(REQUESTID, strconv.FormatUint(id, 10))

		log := logrus.WithFields(logrus.Fields{"logid": id})
		ctx.Set(LOGCTX, log)

		//打印请求日志
		//log.Debugf("req, uri:(%s) body:(%s)", ctx.Request.URL.Path, string(data))

		ctx.Next()
		//打印回复日志
	}
}

func Log(ctx *gin.Context) *logrus.Entry {
	if log, exist := ctx.Get(LOGCTX); exist {
		return log.(*logrus.Entry)
	}
	return logrus.NewEntry(logrus.StandardLogger())
}

// fun(c *gin.Context, req interface{}) (interface{}, error)
// reqfun() (interface{})
// rspfun(c *gin.Context, err error) (interface{})
func Wrapper(fun interface{}, reqfun interface{}, rspfun interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var rsp interface{}

		defer func() {
			if err := recover(); err != nil {
				rspV := reflect.ValueOf(rspfun).Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(err)})
				if len(rspV) == 1 {
					rsp = rspV[0].Interface()
				}
			}

			Log(ctx).Infof("[RSP] URI:(%s) BODY:(%s)", ctx.Request.URL.Path, loghelper.ToPrintString(rsp))

			if _, ok := rsp.(proto.Message); ok {
				ctx.ProtoBuf(http.StatusOK, rsp)
			} else {
				//ctx.Header("Content-Type", "application/json")
				ctx.JSON(http.StatusOK, rsp)
			}
		}()
		reqV := reflect.ValueOf(reqfun).Call([]reflect.Value{})
		if len(reqV) != 1 || reqV[0].IsNil() {
			Log(ctx).Errorf("get req error.")
			panic(errors.New("req struct err"))
		}

		req := reqV[0].Interface()
		var err error
		if _, ok := req.(proto.Message); ok {
			err = ctx.ShouldBindWith(reqV[0].Interface(), binding.ProtoBuf)
		} else {
			err = ctx.ShouldBindJSON(reqV[0].Interface())
		}
		if err != nil {
			Log(ctx).Errorf("parse param err:%v", err)
			panic(err)
		}

		Log(ctx).Infof("[REQ] URI:(%s) BODY:(%s)", ctx.Request.URL.Path, loghelper.ToPrintString(reqV[0].Interface()))

		rspVV := reflect.ValueOf(fun).Call([]reflect.Value{reflect.ValueOf(ctx), reqV[0]})
		if err, ok := rspVV[1].Interface().(error); ok {
			Log(ctx).Errorf("process err: %s", err.Error())
			panic(err)
		}
		rsp = rspVV[0].Interface()
	}
}
