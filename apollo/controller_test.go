package apollo

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_getExchanges(t *testing.T) {
	answer := []string{"a", "b", "c", "d", "e", "f", "g"}

	data := `names=["a", "b", "c", "d", "e", "f", "g"]`

	err := ioutil.WriteFile(tomlFileanme, []byte(data), 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tomlFileanme)

	es := getExchanges()
	for i, v := range answer {
		if es.Names[i] != v {
			t.Error("getExchanges()无法正确读取apollo.toml的内容")
		}
	}
}

// //TODO: 删除此处内容
// func Test_Error(t *testing.T) {
// 	t.Error("看看不通过是什么反应")
// }
