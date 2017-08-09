package exchanges

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	ast := assert.New(t)

	actual := Config("abc")
	expected := "abc.toml"
	ast.Equal(expected, actual, "生成的 exchanges 的配置名称，不对")

}

func Test_Path(t *testing.T) {
	ast := assert.New(t)

	v := url.Values{
		"a": []string{"A"},
		"b": []string{"B"},
		"c": []string{"C"},
	}

	actual := Path("abc", v)
	expected := "abc?a=A&b=B&c=C"
	ast.Equal(expected, actual, "无法合成所需的网址")

}
func Test_MD5(t *testing.T) {
	ast := assert.New(t)

	input := "12345"
	actual := MD5([]byte(input))
	expected := strings.ToLower("827CCB0EEA8A706C4C34A16891F84E7B")
	ast.Equal(expected, HexEncodeToString(actual), "MD5加密后的内容不对。")
}
