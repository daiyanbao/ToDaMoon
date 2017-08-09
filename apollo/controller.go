package apollo

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/aQuaYi/GoKit"
)

var tomlFileanme = "apollo.toml"

//contrller 用来读取tomlFilename的内容
//并根据tomlFilename中的内容，来启动相应的交易所
func controller() {
	ecs := getExchanges()
	for _, e := range ecs.Names {
		fmt.Println(e)
		// switch e {
		// 	case "btc38":

		// 	default:
		// 	log.f
		// }
	}
}

type exchanges struct {
	Names []string
}

func getExchanges() exchanges {
	es := exchanges{}
	if _, err := toml.DecodeFile(tomlFileanme, &es); err != nil {
		msg := fmt.Sprintf("无法加载%s/%s，并Decode到cfg变量: %s", GoKit.PWD(), tomlFileanme, err)
		log.Fatalf(msg)
	}

	return es
}
