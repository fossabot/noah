/*
 * Copyright (c) 2019 Ready Stock
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

package testsuite

import (
    "fmt"
    "github.com/jackc/pgx"
    "github.com/readystock/golog"
)

type TestingLogger interface {
    Log(args ...interface{})
}

type Logger struct {

}

func NewLogger() *Logger {
    return &Logger{}
}

func (l *Logger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
    logArgs := make([]interface{}, 0, 2+len(data))
    logArgs = append(logArgs, level, msg)
    for k, v := range data {
        logArgs = append(logArgs, fmt.Sprintf("%s=%v", k, v))
    }
    switch level {
    case pgx.LogLevelTrace:
        golog.Verbosef("\t[PGX] %v | %s", logArgs, msg)
    case pgx.LogLevelDebug:
        golog.Debugf("\t[PGX] %v | %s", logArgs, msg)
    case pgx.LogLevelInfo:
        golog.Infof("\t[PGX] %v | %s", logArgs, msg)
    case pgx.LogLevelWarn:
        golog.Warnf("\t[PGX] %v | %s", logArgs, msg)
    case pgx.LogLevelError:
        golog.Errorf("\t[PGX] %v | %s", logArgs, msg)
    }
}