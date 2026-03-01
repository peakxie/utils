package ginmw

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testReq struct {
	Name string `json:"name" form:"name"`
}

type testRsp struct {
	Greeting string `json:"greeting"`
}

func testHandler(c *gin.Context, req *testReq) (*testRsp, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	return &testRsp{Greeting: "hello " + req.Name}, nil
}

// bizError implements Coder for testing custom error codes.
type bizError struct {
	code    int
	message string
}

func (e *bizError) Error() string   { return e.message }
func (e *bizError) Code() int       { return e.code }
func (e *bizError) Message() string { return e.message }

func errorHandler(c *gin.Context, req *testReq) (*testRsp, error) {
	return nil, &bizError{code: 1001, message: "user not found"}
}

func plainErrorHandler(c *gin.Context, req *testReq) (*testRsp, error) {
	return nil, errors.New("something went wrong")
}

func parseAPIResponse(t *testing.T, w *httptest.ResponseRecorder) APIResponse {
	t.Helper()
	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v, body: %s", err, w.Body.String())
	}
	return resp
}

func TestWrapperH_Success_POST(t *testing.T) {
	router := gin.New()
	router.POST("/test", WrapperH(testHandler))

	w := httptest.NewRecorder()
	body := `{"name":"world"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	resp := parseAPIResponse(t, w)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("expected message 'success', got %q", resp.Message)
	}

	data, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be map, got %T", resp.Data)
	}
	if data["greeting"] != "hello world" {
		t.Errorf("expected greeting 'hello world', got %v", data["greeting"])
	}
}

func TestWrapperH_Success_GET(t *testing.T) {
	router := gin.New()
	router.GET("/test", WrapperH(testHandler))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test?name=world", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	resp := parseAPIResponse(t, w)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}

	data, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be map, got %T", resp.Data)
	}
	if data["greeting"] != "hello world" {
		t.Errorf("expected greeting 'hello world', got %v", data["greeting"])
	}
}

func TestWrapperH_PlainError(t *testing.T) {
	router := gin.New()
	router.POST("/test", WrapperH(plainErrorHandler))

	w := httptest.NewRecorder()
	body := `{"name":"test"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	resp := parseAPIResponse(t, w)
	if resp.Code != -1 {
		t.Errorf("expected code -1, got %d", resp.Code)
	}
	if resp.Message != "something went wrong" {
		t.Errorf("expected message 'something went wrong', got %q", resp.Message)
	}
	if resp.Data != nil {
		t.Errorf("expected data nil, got %v", resp.Data)
	}
}

func TestWrapperH_CoderError(t *testing.T) {
	router := gin.New()
	router.POST("/test", WrapperH(errorHandler))

	w := httptest.NewRecorder()
	body := `{"name":"test"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	resp := parseAPIResponse(t, w)
	if resp.Code != 1001 {
		t.Errorf("expected code 1001, got %d", resp.Code)
	}
	if resp.Message != "user not found" {
		t.Errorf("expected message 'user not found', got %q", resp.Message)
	}
}

func TestWrapperH_InvalidJSON(t *testing.T) {
	router := gin.New()
	router.POST("/test", WrapperH(testHandler))

	w := httptest.NewRecorder()
	body := `{invalid`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	resp := parseAPIResponse(t, w)
	if resp.Code != -1 {
		t.Errorf("expected code -1, got %d", resp.Code)
	}
	if resp.Message == "" {
		t.Error("expected non-empty error message")
	}
}

func TestWrapperH_HandlerError_EmptyName(t *testing.T) {
	router := gin.New()
	router.POST("/test", WrapperH(testHandler))

	w := httptest.NewRecorder()
	body := `{"name":""}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	resp := parseAPIResponse(t, w)
	if resp.Code != -1 {
		t.Errorf("expected code -1, got %d", resp.Code)
	}
	if resp.Message != "name is required" {
		t.Errorf("expected message 'name is required', got %q", resp.Message)
	}
}

// customResponse is a custom response format for testing ResponsePacker.
type customResponse struct {
	Status  string `json:"status"`
	Payload any    `json:"payload,omitempty"`
	ErrMsg  string `json:"err_msg,omitempty"`
}

// customPacker implements ResponsePacker with a different format.
type customPacker struct{}

func (customPacker) PackSuccess(data any) any {
	return &customResponse{Status: "ok", Payload: data}
}

func (customPacker) PackError(err error) any {
	return &customResponse{Status: "fail", ErrMsg: err.Error()}
}

func TestSetResponsePacker_Custom(t *testing.T) {
	SetResponsePacker(&customPacker{})
	defer SetResponsePacker(nil)

	router := gin.New()
	router.POST("/test", WrapperH(testHandler))

	w := httptest.NewRecorder()
	body := `{"name":"world"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp customResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v, body: %s", err, w.Body.String())
	}
	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Status)
	}
	payload, ok := resp.Payload.(map[string]any)
	if !ok {
		t.Fatalf("expected payload to be map, got %T", resp.Payload)
	}
	if payload["greeting"] != "hello world" {
		t.Errorf("expected greeting 'hello world', got %v", payload["greeting"])
	}
}

func TestSetResponsePacker_CustomError(t *testing.T) {
	SetResponsePacker(&customPacker{})
	defer SetResponsePacker(nil)

	router := gin.New()
	router.POST("/test", WrapperH(plainErrorHandler))

	w := httptest.NewRecorder()
	body := `{"name":"test"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp customResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v, body: %s", err, w.Body.String())
	}
	if resp.Status != "fail" {
		t.Errorf("expected status 'fail', got %q", resp.Status)
	}
	if resp.ErrMsg != "something went wrong" {
		t.Errorf("expected err_msg 'something went wrong', got %q", resp.ErrMsg)
	}
}

func TestSetResponsePacker_ResetToDefault(t *testing.T) {
	SetResponsePacker(&customPacker{})
	SetResponsePacker(nil) // reset

	router := gin.New()
	router.POST("/test", WrapperH(testHandler))

	w := httptest.NewRecorder()
	body := `{"name":"world"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	resp := parseAPIResponse(t, w)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("expected message 'success', got %q", resp.Message)
	}
}

func TestSetLogFunc(t *testing.T) {
	var logCalls int
	SetLogFunc(func(c *gin.Context, format string, args ...any) {
		logCalls++
	})
	defer SetLogFunc(nil)

	router := gin.New()
	router.POST("/test", WrapperH(testHandler))

	// successful request — should log REQ and RSP
	w := httptest.NewRecorder()
	body := `{"name":"world"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if logCalls < 2 {
		t.Errorf("expected at least 2 log calls, got %d", logCalls)
	}

	// error request — should log error
	logCalls = 0
	w = httptest.NewRecorder()
	body = `{"name":""}`
	req, _ = http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if logCalls < 1 {
		t.Errorf("expected at least 1 log call, got %d", logCalls)
	}
}

func TestSetLogFunc_Nil(t *testing.T) {
	// Without LogFunc set, should still work (no panic, no logging).
	SetLogFunc(nil)

	router := gin.New()
	router.POST("/test", WrapperH(testHandler))

	w := httptest.NewRecorder()
	body := `{"name":"world"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	resp := parseAPIResponse(t, w)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
}

func ExampleWrapperH() {
	type LoginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type LoginRsp struct {
		Token string `json:"token"`
	}

	router := gin.Default()
	router.POST("/login", WrapperH(func(c *gin.Context, req *LoginReq) (*LoginRsp, error) {
		if req.Username == "" {
			return nil, fmt.Errorf("username is required")
		}
		return &LoginRsp{Token: "abc123"}, nil
	}))
	// router.Run(":8080")
}
