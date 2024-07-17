package common

import "fmt"

var (
	ErrPassword        = fmt.Errorf("密码不正确")
	ErrUserNameNotNull = fmt.Errorf("用户名不能为空")
)
