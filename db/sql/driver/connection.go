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

package driver

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type Connection interface {
    QueryBytes(sql string, args ...interface{}) ([][]byte, error)
}

type ConnConfig struct {
    Host     string // host (e.g. localhost) or path to unix domain socket directory (e.g. /private/tmp)
    Port     uint16 // default: 5432
    Database string
    User     string
    Password string
}

func (cc *ConnConfig) NetworkAddress() (network, address string) {
    network = "tcp"
    address = fmt.Sprintf("%s:%d", cc.Host, cc.Port)
    // Make sure that the provided address is balid.
    if _, err := os.Stat(cc.Host); err == nil {
        network = "unix"
        address = cc.Host
        if !strings.Contains(address, "/.s.PGSQL.") {
            address = filepath.Join(address, ".s.PGSQL.") + strconv.FormatInt(int64(cc.Port), 10)
        }
    }
    return network, address
}
