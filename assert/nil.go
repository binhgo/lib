package assert

import (
	"github.com/modern-go/test/testify/assert"
	"lib"
	"lib/errors"
	"runtime"
	"strings"
)

func Nil(actual interface{}, errNo int, errMsg string, extMsg ...string) {
	t := CurrentT()
	if assert.Nil(t, actual) {
		return
	}
	Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Error("check failed")
		return
	}
	lib.Log.Warn(ExtractFailedLines(file, line))

	panic(errors.NewUserError(errNo, errMsg, strings.Join(extMsg, ",")))
}

func NotNil(actual interface{}, errNo int, errMsg string, extMsg ...string) {
	t := CurrentT()
	if assert.NotNil(t, actual) {
		return
	}
	Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Error("check failed")
		return
	}
	lib.Log.Warn(ExtractFailedLines(file, line))

	panic(errors.NewUserError(errNo, errMsg, strings.Join(extMsg, ",")))
}
