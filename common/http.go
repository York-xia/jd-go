package common

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	debug *log.Logger
)

func init() {
	debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// http请求接口
type Service interface {
	// 执行Get请求
	Get(url string, args interface{}) ([]byte, error)
}

// http请求默认实现(json传参)
type ServiceImpl struct {
	client *http.Client
}

func (s *ServiceImpl) Get(address string, v interface{}) ([]byte, error) {
	urlVal, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	values, _ := query.Values(v)
	urlVal.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", urlVal.String(), nil)
	if err != nil {
		return nil, err
	}
	a := req.URL.String()
	debug.Println("[京东联盟] Request URI: ", a)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode != 200 {
		return nil, fmt.Errorf("statusCode: %d", statusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func NewService() Service {
	return &ServiceImpl{
		client: &http.Client{},
	}
}
