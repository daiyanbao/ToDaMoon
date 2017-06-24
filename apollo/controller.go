package apollo

import (
	"ToDaMoon/util"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var tomlFileanme = "apollo.toml"

//contrller 用来读取tomlFilename的内容
//并根据tomlFilename中的内容，来启动相应的交易所
func controller() {

}

type exchanges struct {
	Names []string
}

func getExchanges() exchanges {
	es := exchanges{}
	if _, err := toml.DecodeFile(tomlFileanme, &es); err != nil {
		msg := fmt.Sprintf("无法加载%s/%s，并Decode到cfg变量: %s", util.PWD(), tomlFileanme, err)
		log.Fatalf(msg)
	}

	return es
}
