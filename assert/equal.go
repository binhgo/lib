package assert

import (
	"github.com/stretchr/testify/assert"
	"lib"
	"lib/errors"
	"runtime"
	"strings"
)

//go:noinline
func Equal(expected interface{}, actual interface{}, errNo int, errMsg string, extMsg ...string) {
	t := CurrentT()
	if assert.Equal(t, expected, actual) {
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

//go:noinline
func NotEqual(expected interface{}, actual interface{}, errNo int, errMsg string, extMsg ...string) {
	t := CurrentT()
	if assert.NotEqual(t, expected, actual) {
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
