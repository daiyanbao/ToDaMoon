package GoKit

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

//ExternalIP 返回本机的外网IP地址
func ExternalIP() (string, error) {
	addr := "http://myexternalip.com/raw"

	resp, err := http.Get(addr)
	if err != nil {
		return "", Err(err, "无法访问%s", addr)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", Err(err, "无法解析%s返回的数据", addr)
	}

	exip := string(bytes.TrimSpace(b))
	return exip, nil
}
