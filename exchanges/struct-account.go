package exchanges

import (
	"errors"
	"fmt"
	"io/ioutil"

	json "github.com/json-iterator/go"
)

//Account 是用户的账户
type Account struct {
	Coins    map[string]CoinStatus
	TotalCNY float64
}

//CoinStatus 是每个coin在交易所中的状态
//coin也包括cny,usd等真实的货币。
type CoinStatus struct {
	Available float64
	Freezed   float64
}

func readAccountJSON(exchangeName string) (*Account, error) {
	fileName := getAccountFileName(exchangeName)

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		msg := fmt.Sprintf("或者读取%s失败:%s", fileName, err)
		return nil, errors.New(msg)
	}
	a := &Account{}

	if err := json.Unmarshal(bytes, a); err != nil {
		msg := fmt.Sprintf("转换%s文件时失败：%s", fileName, err)
		return nil, errors.New(msg)
	}

	return a, nil
}

//
func saveAccountJSON(exchangeName string, a *Account) error {
	fileName := getAccountFileName(exchangeName)

	bytes, err := json.Marshal(a)
	if err != nil {
		msg := fmt.Sprintf("Marshal%v为JSON时，失败：%s", *a, err)
		return errors.New(msg)
	}

	if err := ioutil.WriteFile(fileName, bytes, 0777); err != nil {
		msg := fmt.Sprintf("保存JSON为%s时，失败：%s", fileName, err)
		return errors.New(msg)
	}
	return nil
}

func getAccountFileName(exchangeName string) string {
	return fmt.Sprintf("%s_account.json", exchangeName)
}
