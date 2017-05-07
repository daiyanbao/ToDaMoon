package util

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

//ExternalIP 返回本机的外网IP地址
func ExternalIP() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", errors.New("无法访问myexternalip.com")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("无法解析myexternalip.com返回的数据")
	}
	exip := string(bytes.TrimSpace(b))
	return exip, nil
}
