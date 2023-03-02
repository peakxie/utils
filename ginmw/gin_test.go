package ginmw_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/peakxie/utils/ginmw"
	"github.com/peakxie/utils/test_proto"
)

type TestReq struct {
	Name    string
	Picture []byte
}

type TestRsp struct {
	RequestId string
	Error     string
}

func GetReq() *TestReq {
	return &TestReq{}
}

func GetRsp(c *gin.Context, err error) *TestRsp {
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

func GetReqPB() *test_proto.Test {
	return &test_proto.Test{}
}

func GetRspPB(c *gin.Context, err error) *test_proto.Test {
	return &test_proto.Test{
		BusinessId: proto.Uint64(3),
		ModuleId:   proto.Uint64(4),
		Picture:    proto.String("ttttttttttttttttttttttttttttttttttttttttttttttttttt"),
		RequestId:  proto.String(err.Error()),
	}
}
func ControlPB(c *gin.Context, req *test_proto.Test) (*test_proto.Test, error) {
	if req.GetRequestId() == "test" {
		return nil, errors.New("name test")
	}
	return &test_proto.Test{BusinessId: req.BusinessId, ModuleId: req.ModuleId, Picture: req.Picture, RequestId: req.RequestId}, nil
}

func TestMain(m *testing.M) {
	//logs.Init("./log", 10, 10, 10, 6)
	router := gin.Default()
	router.POST("/test", ginmw.Wrapper(Control, GetReq, GetRsp))
	router.POST("/testPB", ginmw.Wrapper(ControlPB, GetReqPB, GetRspPB))

	router.Run(":8086")
}
