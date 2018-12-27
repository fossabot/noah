/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

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
