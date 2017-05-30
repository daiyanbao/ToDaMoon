// err.go
// Reference: https://github.com/reusee/codes/tree/master/err
//
//请牢记一点
//error是**接口**
//*Err实现Error()后，是满足error接口的类型
//

package util

import "fmt"

//Error 是自定义的错误类型
type Error struct {
	Info string
	Prev error
}

//不能把Error()改成String()
//因为error是一个接口，实现了Error方法
func (e *Error) Error() string {
	if e.Prev == nil {
		return fmt.Sprintf("%s", e.Info)
	}
	return fmt.Sprintf("%s\n%v", e.Info, e.Prev)
}

//Err 对Error添加新的信息，以便于追踪错误。
func Err(information string, err error) *Error {
	return &Error{
		Info: information,
		Prev: err,
	}
}
