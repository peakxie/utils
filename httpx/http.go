package httpx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/proto"
)

// http post远程调用，使用json协议
func PostJ(client *http.Client, url, contentType string, request interface{}, response interface{}) ([]byte, error) {
	var reqBody []byte
	var err error
	var rsp *http.Response

	if request != nil {
		reqBody, err = json.Marshal(request)
		if err != nil {
			return nil, err
		}
	}

	cType := "application/json"
	if len(contentType) != 0 {
		cType = contentType
	}

	if client != nil {
		rsp, err = client.Post(url, cType, bytes.NewReader(reqBody))
	} else {
		rsp, err = http.Post(url, cType, bytes.NewReader(reqBody))
	}

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

// http post远程调用，使用protobuf协议
func PostP(client *http.Client, url, contentType string, request proto.Message, response proto.Message) ([]byte, error) {
	var reqBody []byte
	var err error
	var rsp *http.Response

	if request != nil {
		reqBody, err = proto.Marshal(request)
		if err != nil {
			return nil, err
		}
	}

	cType := "application/json"
	if len(contentType) != 0 {
		cType = contentType
	}

	if client != nil {
		rsp, err = client.Post(url, cType, bytes.NewReader(reqBody))
	} else {
		rsp, err = http.Post(url, cType, bytes.NewReader(reqBody))
	}

	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody, _ := ioutil.ReadAll(rsp.Body)
	if response != nil {
		err = proto.Unmarshal(rspBody, response)
		if err != nil {
			return rspBody, err
		}
	}
	return rspBody, nil
}
