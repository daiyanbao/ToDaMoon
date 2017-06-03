package exchanges

import (
	"ToDaMoon/util"
	"crypto/md5"
	"encoding/hex"
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

//AskChan 是发送网络请求的通道
type AskChan chan Ask

//Ask 是发送网络请求的内容
type Ask struct {
	Type       askType
	Path       string
	Body       io.Reader
	AnswerChan chan<- Answer
}

type askType int

const (
	//GET method
	GET askType = iota
	//POST method
	POST
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
			case GET:
				data, err := n.get(ask.Path)
				ask.AnswerChan <- Answer{Body: data, Err: err}
			case POST:
				data, err := n.post(ask.Path, ask.Body)
				ask.AnswerChan <- Answer{Body: data, Err: err}
			default:
				log.Println("错误的请求类型")
			}
			beginTime = util.HoldOn(waitTime, beginTime)
		}
	}()
}
func (n *Net) ask(t askType, path string, body io.Reader) ([]byte, error) {
	ansChan := make(chan Answer)
	n.Ask <- Ask{Type: t,
		Path:       path,
		Body:       body,
		AnswerChan: ansChan,
	}

	ans := <-ansChan
	if ans.Err != nil {
		return nil, ans.Err
	}

	return ans.Body, nil
}

//Post 是一个交易所网络核心的通用Post方法。
func (n *Net) Post(path string, body io.Reader) ([]byte, error) {
	return n.ask(POST, path, body)
}

//Get 是一个交易所网络核心的通用Get方法。
func (n *Net) Get(path string) ([]byte, error) {
	return n.ask(GET, path, nil)
}

//Path 制造请求的网址
func Path(URL string, values url.Values) string {
	path := URL
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

//MD5 对input的内容进行md5加密
func MD5(input []byte) []byte {
	hash := md5.New()
	hash.Write(input)
	return hash.Sum(nil)
}

//HexEncodeToString 对input的内容进行hex加密
func HexEncodeToString(input []byte) string {
	return hex.EncodeToString(input)
}
