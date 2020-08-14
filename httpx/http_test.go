package httpx_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peakxie/utils/httpx"
)

//test http接口

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
