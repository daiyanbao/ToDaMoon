package Interface

//Notify 是通知类型的接口
type Notify interface {
	//反馈程序运行状态。
	Debug(string) error
	Warning(string) error
	Error(string) error

	//反馈系统运行信息
	Info(string) error
	Good(string) error
	Bad(string) error
}
