//Package pubu API的简单封装，并根据个人需求，设置了一些预定义的消息样式。
package pubu

import "net/http"
import "io/ioutil"
import "encoding/json"
import "log"

//Client 是发送消息的客户端。
type Client struct {
	hook string
}

//Send 给服务器发送消息
func (c *Client) Send(ic *Incoming) error {
	message, err := ic.Build()
	if err != nil {
		return err
	}

	req := &http.Client{}
	resp, err := req.Post(c.hook, "application/json", message)

	if err != nil {
		return err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type r struct {
		Err  int `json:"error"`
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	var ans r

	json.Unmarshal(result, &ans)
	if ans.Err != 0 {
		log.Println("发送消息出错。", ic, "出错原因：", ans.Data.Message)
	}

	resp.Body.Close()
	return nil
}

//New 创建新的客户端。
func New(hook string) *Client {
	return &Client{hook: hook}
}
