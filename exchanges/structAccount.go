package exchanges

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
)

//Account 是用户的账户
type Account struct {
	Coins    map[string]CoinStatus
	TotalCNY float64
}

//NewAccount 返回*Account
func NewAccount() *Account {
	return &Account{
		Coins: make(map[string]CoinStatus),
	}
}

//CoinStatus 是每个coin在交易所中的状态
//coin也包括cny,usd等真实的货币。
type CoinStatus struct {
	Available float64
	Freezed   float64
	Total     float64
}

func (cs CoinStatus) String() string {
	return fmt.Sprintf("\t%f\t%f\t%f\n", cs.Available, cs.Freezed, cs.Total)
}

func (a Account) String() string {
	result := ""

	result += fmt.Sprintf("总金额: ￥%.2f\n", a.TotalCNY)
	result += fmt.Sprintln("Coin\t可用\t\t冻结\t\t总计")

	sortedCS := sortedCoins(a.Coins)

	for _, k := range sortedCS {
		result += k + a.Coins[k].String()
	}

	return result
}

func sortedCoins(m map[string]CoinStatus) []string {
	result := make([]string, len(m))
	i := 0

	for k := range m {
		result[i] = k
		i++
	}

	sort.Strings(result)
	return result
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
