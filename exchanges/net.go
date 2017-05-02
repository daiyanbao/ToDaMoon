//Package exchanges 的net.go提供了交易所API访问所需的方法。
package exchanges

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//Get 网络数据获取方式的封装
func Get(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("HTTP status code: %d\n", res.StatusCode)
		return nil, errors.New("Status code was not 200.")
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

//Post 网络数据获取方式的封装
func Post(path string,
	headers map[string]string,
	body io.Reader,
) ([]byte, error) {
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

//Path 制造请求的网址
func Path(url string, values url.Values) string {
	path := url
	if len(values) > 0 {
		path += "?" + values.Encode()
	}
	return path
}

//JSONEncode 把数据转换成json格式
func JSONEncode(v interface{}) ([]byte, error) {
	json, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}

	return json, nil
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
