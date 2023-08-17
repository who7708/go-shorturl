package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	// 导入当前包，fix main.App not found
	"app"
)

const (
	expTime       = 60
	longURL       = "https://www.baidu.com"
	shortUrl      = "Hello"
	shortLinkInfo = `{ "url":"https://www.baidu.com", "create_at":"2023-08-17 21:17:12.34971 +0800 CST m=+62.079104201", "expiration_in_minutes":1 }`
)

type storageMock struct {
	mock.Mock
}

var app main.App
var mockR *storageMock

func (s *storageMock) Shorten(url string, expirationInMinutes int64) (string, error) {
	args := s.Called(url, expirationInMinutes)
	return args.String(0), args.Error(1)
}

func (s *storageMock) ShortlinkInfo(shortUrl string) (interface{}, error) {
	args := s.Called(shortUrl)
	return args.String(0), args.Error(1)
}

func (s *storageMock) Unshortend(shortUrl string) (string, error) {
	args := s.Called(shortUrl)
	return args.String(0), args.Error(1)
}

func init() {
	app = main.App{}
	mockR = new(storageMock)
	app.Initialize(&main.Env{S: mockR})
}

func TestCrateShortLink(t *testing.T) {
	var jsonStr = []byte(`{
    "url": "https://www.baidu.com",
    "expiration_in_minutes": 60
}`)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal("无法创建请求", err)
	}
	req.Header.Set("Content-TYpe", "application/json")

	// mock
	mockR.On("Shorten", longURL, int64(expTime)).Return(shortUrl, nil).Once()

	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatal("期望值：%d, 实际: %d", http.StatusCreated, rw.Code)
	}

	resp := struct {
		ShortLink string `json:"shortlink"`
	}{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatal("解码错误")
	}

	if resp.ShortLink != shortUrl {
		t.Fatal("期望值：%s, 实际: %s", shortUrl, resp.ShortLink)
	}

}

func TestShortLinkRedirect(t *testing.T) {
	r := fmt.Sprintf("/%s", shortUrl)

	req, err := http.NewRequest("GET", r, nil)

	if err != nil {
		t.Fatal("无法创建请求", err)
	}

	mockR.On("Unshortend", shortUrl).Return(longURL, nil).Once()

	rw := httptest.NewRecorder()

	app.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusTemporaryRedirect {
		t.Fatal("期望值：%s, 实际: %s", http.StatusTemporaryRedirect, rw.Code)
	}
}
