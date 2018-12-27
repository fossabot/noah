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

package util

import (
    "fmt"
    "github.com/kataras/go-errors"
    "github.com/kataras/golog"
)

func CatchPanic(err *error) {
    if x := recover(); x != nil {
        golog.Errorf("uncaught panic: %v", x)
        *err = fmt.Errorf("uncaught panic: %v", x)
    }
}

func CombineErrors(errs []error) error {
    if len(errs) > 0 {
        err := errors.New(errs[0].Error())
        for i := 1; i < len(errs); i++ {
            err.AppendErr(errs[i])
        }
        return err
    }
    return nil
}
