package exchanges

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Net_Get(t *testing.T) {
	ast := assert.New(t)
	msg := "hello, client"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, msg)
	}))
	defer ts.Close()

	n := NewNet(100)

	data, err := n.Get(ts.URL)
	ast.Nil(err, "n.Get，出错了")
	ast.Equal(msg+"\n", string(data), "ts发来的数据不对")
}

func Test_Net_Get_noExistServer(t *testing.T) {
	ast := assert.New(t)
	n := NewNet(100)

	_, err := n.Get("localhost:8888")
	ast.NotNil(err, "n.Get访问不存在的服务器，却没有报错")
}

func Test_Net_sleep(t *testing.T) {
	ast := assert.New(t)
	msg := "hello, client"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, msg)
	}))
	defer ts.Close()

	sleepMS := 50
	num := 10
	n := NewNet(sleepMS)

	begin := time.Now()
	for i := 0; i < num; i++ {
		n.Get(ts.URL)
	}
	ast.WithinDuration(begin.Add(time.Millisecond*time.Duration(sleepMS*num)), time.Now(), time.Millisecond*time.Duration(sleepMS), "Net没有按照规划的时间休息")
}

func Test_Net_Post(t *testing.T) {
	ast := assert.New(t)
	msg := "hello, post"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, msg)
	}))
	defer ts.Close()

	n := NewNet(100)

	data, err := n.Post(ts.URL, nil)
	ast.Nil(err, "n.Post，出错了")
	ast.Equal(msg+"\n", string(data), "ts发来的数据不对")
}

func Test_Net_post_wrongURL(t *testing.T) {
	ast := assert.New(t)

	n := NewNet(100)

	_, err := n.Post("::::::::::", nil)
	ast.NotNil(err, "n.Post使用了不存在的地址，却没有报错")
}

func Test_Net_post_httpClientDo(t *testing.T) {
	ast := assert.New(t)
	n := NewNet(100)

	_, err := n.Post("", nil)
	ast.NotNil(err, "n.Post的body为nil，却没有报错")
}

func Test_handleResp(t *testing.T) {
	ast := assert.New(t)
	r := new(http.Response)

	r.StatusCode = 404
	r.Body = ioutil.NopCloser(strings.NewReader(""))

	resp := handleResp(r)
	ast.NotNil(resp.err, "handleResp 没能识别出http.Response中的状态码不为200")
}
