package httpx_test

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/peakxie/utils/httpx"
	"github.com/peakxie/utils/test_proto"
)

//test http接口
/*
type JsonTest struct {
	Name    string
	Picture string
}

func server() {

	go func() {
		e := gin.Default()
		e.POST("/test", func(c *gin.Context) {
			req := &JsonTest{}
			if e := c.ShouldBindJSON(req); e != nil {
				fmt.Println("gin bind req", e)
				return
			}
			req.Name += "_rsp"
			c.JSON(http.StatusOK, req)
		})

		e.Run(":8080")
	}()
}

func TestHttpJ(t *testing.T) {
	server()
	for {
		time.Sleep(time.Second * 3)

		var rsp JsonTest
		//发请求，自动打印请求包和回包，不需要自己序列化和反序列化协议
		_, e := httpx.PostJ(nil, "http://localhost:8080/test", "", &JsonTest{Name: "name", Picture: "xxxxxxxxxxxxxxx"}, &rsp)
		if e != nil {
			fmt.Println("err", e)
		} else {
			fmt.Println("succ: ", rsp)
		}
	}
}
*/

type TestReq struct {
	Name    string
	Picture []byte
}

type TestRsp struct {
	RequestId string
	Error     string
}

func TestHttp(t *testing.T) {
	//发请求，自动打印请求包和回包，不需要自己序列化和反序列化协议
	req := &TestReq{
		Name:    "test",
		Picture: []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
	}
	rsp := &TestRsp{}
	_, e := httpx.PostJ(nil, "http://localhost:8086/test", "", req, rsp)
	if e != nil {
		fmt.Println("err", e)
	} else {
		fmt.Printf("succ: %+v\n", rsp)
	}

	req.Name = "test1"
	_, e = httpx.PostJ(nil, "http://localhost:8086/test", "", req, rsp)
	if e != nil {
		fmt.Println("err", e)
	} else {
		fmt.Printf("succ: %+v\n", rsp)
	}

	reqPB := &test_proto.Test{
		BusinessId: proto.Uint64(1),
		ModuleId:   proto.Uint64(2),
		Picture:    proto.String("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		RequestId:  proto.String("test"),
	}
	rspPB := &test_proto.Test{}
	_, e = httpx.PostP(nil, "http://localhost:8086/testPB", "", reqPB, rspPB)
	if e != nil {
		fmt.Println("err", e)
	} else {
		fmt.Println("succ: ", rspPB.String())
	}

	reqPB.RequestId = proto.String("test1")
	_, e = httpx.PostP(nil, "http://localhost:8086/testPB", "", reqPB, rspPB)
	if e != nil {
		fmt.Println("err", e)
	} else {
		fmt.Println("succ: ", rspPB.String())
	}

}
