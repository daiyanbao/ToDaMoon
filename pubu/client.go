package pubu

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type incomingChan chan *incoming

type client struct {
	hook   string
	icChan incomingChan
}

func (c *client) send(ic *incoming) {
	message, err := ic.build()
	if err != nil {
		log.Println("发送消息出错。", ic, "出错原因：", err)
		return
	}

	req := &http.Client{}
	resp, err := req.Post(c.hook, "application/json", message)

	if err != nil {
		log.Println("发送消息出错。", ic, "出错原因：", err)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("发送消息出错。", ic, "出错原因：", err)
		return
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
	return
}

func (c *client) Debug(msg string) {
	m := msgMaker(msg)
	m.debug()
	c.icChan <- m
}

func (c *client) Warning(msg string) {
	m := msgMaker(msg)
	m.warning()
	c.icChan <- m
}

func (c *client) Error(msg string) {
	m := msgMaker(msg)
	m.error()
	c.icChan <- m
}

func (c *client) Info(msg string) {
	m := msgMaker(msg)
	m.info()
	c.icChan <- m
}

func (c *client) Good(msg string) {
	m := msgMaker(msg)
	m.good()
	c.icChan <- m
}

func (c *client) Bad(msg string) {
	m := msgMaker(msg)
	m.bad()
	c.icChan <- m
}
