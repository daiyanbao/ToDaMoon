// net.go
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

func Get(url string) (contents []byte, err error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Printf("HTTP status code: %d\n", res.StatusCode)
		return nil, errors.New("Status code was not 200.")
	}

	contents, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	res.Body.Close()

	return contents, nil
}

func Post(path string,
	headers map[string]string,
	body io.Reader,
) (contents []byte, err error) {
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

	contents, err = ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return contents, nil
}

func Path(url string, values url.Values) string {
	path := url
	if len(values) > 0 {
		path += "?" + values.Encode()
	}
	return path
}

func JSONEncode(v interface{}) ([]byte, error) {
	json, err := json.Marshal(&v)

	if err != nil {
		return nil, err
	}

	return json, nil
}

func JSONDecode(data []byte, to interface{}) error {
	err := json.Unmarshal(data, &to)

	if err != nil {
		return err
	}

	return nil
}

func MD5(input []byte) []byte {
	hash := md5.New()
	hash.Write(input)
	return hash.Sum(nil)
}

func HexEncodeToString(input []byte) string {
	return hex.EncodeToString(input)
}
