package exchanges

import (
	"fmt"
	"log"
)

//这个部分是抽象了exchagnes管理trades数据库方法的内容

//TradesDBs 是exchange中存储coin历史交易记录的数据库
type TradesDBs map[string]map[string]*TradesDB

//MakeTradesDBs 是链接了exchanges所有的Trades数据库
func MakeTradesDBs(dir, exchange string, markets map[string][]string) TradesDBs {
	log.Printf("开始链接%s的数据库", exchange)
	var err error
	t := make(map[string]map[string]*TradesDB)
	for money, coins := range markets {
		t[money] = make(map[string]*TradesDB)
		for _, coin := range coins {
			filename := tradesDBFileName(dir, exchange, money, coin)
			t[money][coin], err = OpenTradesDB(filename)
			if err != nil {
				text := fmt.Sprintf("无法连接%s数据库: %s", filename, err)
				log.Fatalln(text)
			}
			log.Printf("已连接上：%s的%s的%s的数据库%s", exchange, money, coin, filename)
		}
	}
	return t
}

//生成数据库文件路径
func tradesDBFileName(dir, exchange, money, coin string) string {
	return fmt.Sprintf("%s/%s/%s-%s.db", dir, exchange, money, coin)
}
