package errors

import (
	"fmt"
	"strconv"
	"strings"
)

func NewUserError(code int, desc string, ext ...string) error {
	zhDesc := ""
	if len(ext) != 0 {
		zhDesc = strings.Join(ext, ",")
	}
	return &UserError{code, desc, zhDesc}
}

type UserError struct {
	Code   int
	Desc   string
	ZhDesc string
}

func (e *UserError) Error() string {
	return fmt.Sprintf("[%d]%s", e.Code, e.Desc)
}

func (e *UserError) ErrNo() int {
	return e.Code
}

func (e *UserError) ErrNoStr() string {
	return strconv.Itoa(e.Code)
}

func (e *UserError) GetZhDesc() string {
	return fmt.Sprintf("%s", e.ZhDesc)
}
