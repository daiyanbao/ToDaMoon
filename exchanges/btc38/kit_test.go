package btc38

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/json-iterator/go/assert"
)

func Test_price2Str(t *testing.T) {
	data := map[float64]string{
		123456789.12345: "123456789",
		12345.6789:      "12345.6",
		123.456789:      "123.456",
		1.23456789:      "1.23457",
		0.00123456789:   "0.00123",
	}

	for k, v := range data {
		psk := price2Str(k)
		if psk != v {
			t.Errorf("%f应该被转换成%s，而不是%s", k, v, psk)
		}
	}
}
func Test_ExternalIP(t *testing.T) {
	cmd := exec.Command("curl", "http://myexternalip.com/raw")
	out, err := cmd.Output()
	assert.Nil(t, err, "从命令行获取外部IP的方式，失败: %s", err)
	expectedIP := string(bytes.TrimSpace(out))

	actualIP, err := ExternalIP()
	assert.Nil(t, err, "从ExternalIP获取外部IP的方式，失败: %s", err)

	assert.Equal(t, expectedIP, actualIP, "两种方式获取的外部IP不一致。")
}
