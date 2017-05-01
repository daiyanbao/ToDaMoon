package Interface

//Notify 是通知类型的接口
type Notify interface {
	Debug(string) error
	Warning(string) error
	Error(string) error
	Info(string) error
	Good(string) error
	Bad(string) error
}
