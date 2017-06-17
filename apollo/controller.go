package apollo

import (
	"ToDaMoon/util"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var tomlFileanme = "apollo.toml"

func control() {

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
