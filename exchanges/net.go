package exchanges

import (
	"ToDaMoon/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type AskChan chan Ask

type Ask struct {
	Type       askType
	Path       string
	Body       io.Reader
	AnswerChan chan<- Answer
}

type askType int

const (
	//Get method
	Get askType = iota
	//Post method
	Post
)

//Answer 是统一的应答格式
type Answer struct {
	Body []byte
	Err  error
}

//Net 包含了交易所模块网络访问的基础结构
type Net struct {
	Header map[string]string
	Ask    AskChan
}

//Post 网络数据获取方式的封装
func (n *Net) post(path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	//设置访问的请求头
	for k, v := range n.Header {
		req.Header.Add(k, v)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return handleResp(resp)
}

//Get 网络数据获取方式的封装
func (n *Net) get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return handleResp(resp)
}

func handleResp(r *http.Response) ([]byte, error) {
	if r.StatusCode/100 != 2 {
		text := fmt.Sprintf("响应码是%d\n", r.StatusCode)
		return nil, errors.New(text)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

//Start 启动了网络的核心部分
func (n *Net) Start(waitMS int) {
	askCh := make(AskChan, 12)
	n.Ask = askCh
	waitTime := time.Millisecond * time.Duration(waitMS)

	go func() {
		beginTime := time.Now()
		for ask := range n.Ask {
			switch ask.Type {
			case Get:
				data, err := n.get(ask.Path)
				ask.AnswerChan <- Answer{Body: data, Err: err}
			case Post:
				data, err := n.post(ask.Path, ask.Body)
				ask.AnswerChan <- Answer{Body: data, Err: err}
			default:
				log.Println("错误的请求类型")
			}
			beginTime = util.HoldOn(waitTime, beginTime)
		}
	}()
}

//Path 制造请求的网址
func Path(url string, values url.Values) string {
	path := url
	if len(values) > 0 {
		path += "?" + values.Encode()
	}
	return path
}

//JSONDecode 解析Json格式
func JSONDecode(data []byte, to interface{}) error {
	err := json.Unmarshal(data, &to)
	if err != nil {
		return err
	}
	return nil
}
