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

package npgx

import (
    "database/sql/driver"
    "reflect"
)

var valuerReflectType = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

// callValuerValue returns vr.Value(), with one exception:
// If vr.Value is an auto-generated method on a pointer type and the
// pointer is nil, it would panic at runtime in the panicwrap
// method. Treat it like nil instead.
// Issue 8415.
//
// This is so people can implement driver.Value on value types and
// still use nil pointers to those types to mean nil/NULL, just like
// string/*string.
//
// This function is mirrored in the database/sql/driver package.
func callValuerValue(vr driver.Valuer) (v driver.Value, err error) {
    if rv := reflect.ValueOf(vr); rv.Kind() == reflect.Ptr &&
        rv.IsNil() &&
        rv.Type().Elem().Implements(valuerReflectType) {
        return nil, nil
    }
    return vr.Value()
}
