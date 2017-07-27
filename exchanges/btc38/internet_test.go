package GoKit

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExternalIP(t *testing.T) {
	cmd := exec.Command("curl", "http://myexternalip.com/raw") /// 查看当前目录下文件
	out, err := cmd.Output()
	assert.Nil(t, err, "从命令行获取外部IP的方式，失败: %s", err)
	expectedIP := string(bytes.TrimSpace(out))

	actualIP, err := ExternalIP()
	assert.Nil(t, err, "从ExternalIP获取外部IP的方式，失败: %s", err)

	assert.Equal(t, expectedIP, actualIP, "两种方式获取的外部IP不一致。")
}
