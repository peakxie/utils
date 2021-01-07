package ginhelper_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/peakxie/utils/ginhelper"
)

type TestReq struct {
	Name    string
	Picture string
}

type TestRsp struct {
	RequestId string
	Error     string
}

func GetReq() *TestReq {
	return &TestReq{}
}

func GetRsp(err error) *TestRsp {
	return &TestRsp{
		Error: err.Error(),
	}
}

func Control(c *gin.Context, req *TestReq) (*TestRsp, error) {
	if req.Name == "test" {
		return nil, errors.New("name test")
	}
	return &TestRsp{RequestId: req.Name}, nil
}

func TestMain(m *testing.M) {
	//logs.Init("./log", 10, 10, 10, 6)
	router := gin.Default()
	router.POST("/test", ginhelper.WrapperJ(Control, GetReq, GetRsp))

	router.Run(":8086")
}
