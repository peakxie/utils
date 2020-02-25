package httpsw

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var client = &http.Client{}

func Init(c *http.Client) {
	client = c
}

// http post远程调用，使用json协议
func PostJ(ctx context.Context, url string, request interface{}, response interface{}) ([]byte, error) {
	var reqBody []byte
	var err error

	if request != nil {
		reqBody, err = json.Marshal(request)
		if err != nil {
			return nil, err
		}
	}

	rsp, err := client.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody, _ := ioutil.ReadAll(rsp.Body)
	if response != nil {
		err = json.Unmarshal(rspBody, response)
		if err != nil {
			return rspBody, err
		}
	}

	return rspBody, nil
}
