package utility

import (
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestIRGet(t *testing.T) {
	ir := IRRequest{
		logger: &zap.SugaredLogger{},
		path:   "/v1/api/alert/count",
		domain: "http://127.0.0.1:8081",
		params: nil,
	}

	resp, err := ir.Get()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.TypeOf(resp))
}

func TestIRPost(t *testing.T) {
	logger := zap.NewExample().Sugar()
	ir := IRRequest{
		logger: logger,
		path:   "/v1/api/article/publish",
		domain: "http://127.0.0.1:8081/",
		params: map[string]string{"id": "1", "title": "标题1111", "content": "内容内容"},
	}

	resp, err := ir.Post()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err)
	fmt.Println(reflect.TypeOf(resp))
}

func TestRetry(t *testing.T) {
	resp, err := Retry(3, "https://boluome.com/api/filebeat/ping", nil, IRPost)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
	fmt.Println(resp)
}

func TestDecorateRetry(t *testing.T) {
	IRPost := DecorateRetry(5, "https://boluome.tech/api/filebeat/ping", nil, IRPost)
	resp, err := IRPost("https://boluome.tech/api/filebeat/ping", nil)
	fmt.Println(err)
	fmt.Println(resp)
}
