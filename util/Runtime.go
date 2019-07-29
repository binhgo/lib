package util

import (
	"runtime"
)

func GetStackTrace(r interface{}) string{
	//fmt.Printf("Internal error: %v", r)
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)

	return string(buf[0:stackSize])
}
