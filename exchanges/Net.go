package exchanges

/*
因为交易所都会有RESTful API访问频率限制。
所以所有的交易所访问请求，都会由专门的goroutine负责，以便控制访问频率。
*/

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	gk "github.com/aQuaYi/GoKit"
)

// Net 交易所 RESTful API 统一访问接口
type Net interface {
	Post(path string, body io.Reader) ([]byte, error)
	Get(path string) ([]byte, error)
}

// reqChan 是发送API访问请求的通道
type reqChan chan request

// 包含了API访问请求的细节
type request struct {
	method   reqMethod
	path     string
	body     io.Reader
	respChan chan<- *response
}

type reqMethod int

const (
	get reqMethod = iota
	post
)

// 统一的应答格式
type response struct {
	body []byte
	err  error
}

// net 包含了交易所模块网络访问的基础结构
type net struct {
	reqChan reqChan
	header  map[string]string
	sleep   func()
}

// NewNet 返回一个可以立即使用的交易所网络访问核心
func NewNet(sleepMS int) Net {
	n := &net{
		reqChan: make(reqChan, 3),
		header: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   "Mozilla/4.0"},
		sleep: gk.SleepFunc(time.Millisecond * time.Duration(sleepMS)),
	}

	go start(n)

	return n
}

// start 启动了网络的核心部分
func start(n *net) {
	for req := range n.reqChan {
		switch req.method {
		case post:
			req.respChan <- n.post(req.path, req.body)
		default:
			req.respChan <- n.get(req.path)
		}

		n.sleep()
	}
}

// Get 是一个交易所网络核心的通用Get方法。
func (n *net) Get(path string) ([]byte, error) {
	return n.request(get, path, nil)
}

// Post 是一个交易所网络核心的通用Post方法。
func (n *net) Post(path string, body io.Reader) ([]byte, error) {
	return n.request(post, path, body)
}

func (n *net) request(t reqMethod, path string, body io.Reader) ([]byte, error) {
	respChan := make(chan *response)
	n.reqChan <- request{method: t,
		path:     path,
		body:     body,
		respChan: respChan,
	}

	ans := <-respChan

	return ans.body, ans.err
}

//  get 网络数据获取方式的封装
func (n *net) get(url string) *response {
	resp, err := http.Get(url)
	if err != nil {
		return &response{nil, err}
	}

	return handleResp(resp)
}

//  post 网络数据获取方式的封装
func (n *net) post(path string, body io.Reader) *response {
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return &response{nil, err}
	}

	//  设置访问的请求头
	for k, v := range n.header {
		req.Header.Add(k, v)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return &response{nil, err}
	}

	return handleResp(resp)
}

func handleResp(r *http.Response) *response {
	defer r.Body.Close()
	if r.StatusCode/100 != 2 {
		text := fmt.Sprintf("响应码是%d\n", r.StatusCode)
		return &response{nil, errors.New(text)}
	}

	body, err := ioutil.ReadAll(r.Body)

	return &response{body, err}
}
