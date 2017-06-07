package exchanges

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
