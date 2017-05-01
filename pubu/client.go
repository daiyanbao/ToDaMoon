package pubu

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)



func (p *pubu) send(ic *Incoming) error {
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
