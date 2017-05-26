package exchanges

import (
	"ToDaMoon/util"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type askChan chan ask

type ask struct {
	Type       askType
	Path       string
	Body       io.Reader
	AnswerChan chan<- answer
}

type askType int

const (
	get askType = iota
	post
)

type answer struct {
	body []byte
	err  error
}

type net struct {
	Header http.Header
	Ask    askChan
}

//Post 网络数据获取方式的封装
func (n *net) post(path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	//设置访问的请求头
	req.Header = n.Header

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return handleResp(resp)
}

//Get 网络数据获取方式的封装
func (n *net) get(url string) ([]byte, error) {
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

func (n *net) start(waitTime time.Duration) {
	askCh := make(askChan, 12)
	n.Ask = askCh

	go func() {
		beginTime := time.Now()
		for ask := range n.Ask {
			switch ask.Type {
			case get:
				data, err := n.get(ask.Path)
				ask.AnswerChan <- answer{body: data, err: err}
			case post:
				data, err := n.post(ask.Path, ask.Body)
				ask.AnswerChan <- answer{body: data, err: err}
			default:
				log.Println("错误的请求类型")
			}
			beginTime = util.HoldOn(waitTime, beginTime)
		}
	}()
}
