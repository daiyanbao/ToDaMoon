package exchanges

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
)

// Config 根据name返回exchange的配置文件名称
func Config(name string) string {
	return fmt.Sprintf("%s.toml", name)
}

// Path 制造请求的网址
func Path(URL string, values url.Values) string {
	path := URL
	if len(values) > 0 {
		path += "?" + values.Encode()
	}
	return path
}

// MD5 对input的内容进行md5加密
// TODO: 了解一下，MD函数中，发生了什么
func MD5(input []byte) []byte {
	hash := md5.New()
	hash.Write(input)
	return hash.Sum(nil)
}

// HexEncodeToString 对input的内容进行hex加密
// TODO: 到底是怎么用的，要不要和MD5合并
func HexEncodeToString(input []byte) string {
	return hex.EncodeToString(input)
}
