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

package npgx_test

import (
    "github.com/readystock/noah/db/sql/driver"
    "github.com/readystock/noah/db/sql/driver/npgx"
    "testing"
)

var (
    ConnectionConfig = driver.ConnConfig{
        Host:     "localhost",
        Port:     5432,
        User:     "postgres",
        Password: "Spring!2016",
        Database: "postgres",
    }
    ConnectionConfig2 = driver.ConnConfig{
        Host:     "localhost",
        Port:     5432,
        User:     "postgres",
        Password: "Spring!2016",
        Database: "ready_three",
    }

    BadConnectionConfig = driver.ConnConfig{
        Host:     "localhost",
        Port:     123,
        User:     "postgres",
        Password: "Spring!2016",
        Database: "postgres",
    }
)

func Test_Connect(t *testing.T) {
    d, err := npgx.Connect(ConnectionConfig)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    r, err := d.Query("SELECT 1;")
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    for r.Next() {
        fields, err := r.Values()
        if err != nil {
            t.Error(err)
            t.Fail()
        }
        if len(fields) == 1 {
            if fields[0].(int32) != 1 {
                t.Error("Unexpected result!")
                t.Fail()
            }
        }
    }
    err = d.Close()
    if err != nil {
        t.Error(err)
        t.Fail()
    }
}

func Test_CITEXT(t *testing.T) {
    d, err := npgx.Connect(ConnectionConfig)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    r, err := d.Query("SELECT 123::CITEXT;")
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    for r.Next() {
        fields, err := r.Values()
        if err != nil {
            t.Error(err)
            t.Fail()
        }
        if len(fields) == 1 {
            if fields[0].(string) != "123" {
                t.Error("Unexpected result!")
                t.Fail()
            }
        }
    }
    err = d.Close()
    if err != nil {
        t.Error(err)
        t.Fail()
    }
}

func Test_BadConnection(t *testing.T) {
    _, err := npgx.Connect(BadConnectionConfig)
    if err == nil {
        t.Fail()
    }
}
