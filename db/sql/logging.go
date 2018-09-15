package sql

import (
	"fmt"
	"github.com/kataras/golog"
)

func (ex *connExecutor) Debug(msg string, args ...interface{}) {
	golog.Debugf(fmt.Sprintf("[%s] %s", ex.ClientAddress, msg), args...)
}

func (ex *connExecutor) Info(msg string, args ...interface{}) {
	golog.Infof(fmt.Sprintf("[%s] %s", ex.ClientAddress, msg), args...)
}

func (ex *connExecutor) Warn(msg string, args ...interface{}) {
	golog.Warnf(fmt.Sprintf("[%s] %s", ex.ClientAddress, msg), args...)
}

func (ex *connExecutor) Error(msg string, args ...interface{}) {
	golog.Errorf(fmt.Sprintf("[%s] %s", ex.ClientAddress, msg), args...)
}

func (ex *connExecutor) Fatal(msg string, args ...interface{}) {
	golog.Fatalf(fmt.Sprintf("[%s] %s", ex.ClientAddress, msg), args...)
}