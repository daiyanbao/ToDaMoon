package btc38

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aQuaYi/GoKit"
)

func price2Str(price float64) string {
	ps := fmt.Sprintf("%.5f", price)

	i := strings.Index(ps, ".")
	if i == 1 {
		return ps
	}

	remain := 7 - i
	if remain > 0 {
		return ps[:7]
	}

	return ps[:i]
}

// ExternalIP 是获取外部IP的方法
func ExternalIP() (string, error) {
	addr := "http://myexternalip.com/raw"

	resp, err := http.Get(addr)
	if err != nil {
		return "", GoKit.Err(err, "无法访问%s", addr)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", GoKit.Err(err, "无法解析%s返回的数据", addr)
	}

	exip := string(bytes.TrimSpace(b))
	return exip, nil
}
