package apollo

//Notifier 是通知类型的接口
type Notifier interface {
	//反馈程序运行状态。
	Debug(string)
	Warning(string)
	Error(string)

	//反馈系统运行信息
	Info(string)
	Good(string)
	Bad(string)
}
