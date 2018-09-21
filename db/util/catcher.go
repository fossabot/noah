package util

import (
	"fmt"
	"github.com/kataras/golog"
)

func CatchPanic(err *error) {
	if x := recover(); x != nil {
		golog.Errorf("uncaught panic: %v", x)
		*err = fmt.Errorf("uncaught panic: %v", x)
	}
}