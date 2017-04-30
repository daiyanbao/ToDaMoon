//Package pubu API的简单封装，并根据个人需求，设置了一些预定义的消息样式。
package pubu

import "net/http"

//Client 是发送消息的客户端。
type Client struct {
	hook string
	name string
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

	resp.Body.Close()
	return nil
}

//New 创建新的客户端。
func New(name, hook string) *Client {
	return &Client{hook: hook, name: name}
}
