package middleware

import (
	"launchpad/enums"
)

// BusinessError 自定义业务异常
type BusinessError struct {
	ReCode enums.ReCode
}

func (e *BusinessError) Error() string {
	return e.ReCode.GetDesc()
}
