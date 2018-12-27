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

package snowflake

import (
    "fmt"
    "github.com/deckarep/golang-set"
    "testing"
)

var sf *Snowflake

func init() {
    sf = NewSnowflake(1)
    if sf == nil {
        panic("snowflake not created")
    }
}

func nextID(t *testing.T) uint64 {
    id, err := sf.NextID()
    if err != nil {
        t.Fatal("id not generated")
        t.Fail()
    }
    return id
}

func Test_SnowflakeParallel(t *testing.T) {
    const numID = 10000
    const numGenerator = 8
    consumer := make(chan uint64, numID*numGenerator)

    generate := func() {
        for i := 0; i < numID; i++ {
            consumer <- nextID(t)
        }
    }

    for i := 0; i < numGenerator; i++ {
        go generate()
    }

    set := mapset.NewSet()
    for i := 0; i < numID*numGenerator; i++ {
        id := <-consumer
        if set.Contains(id) {
            t.Fatal("duplicated id")
        } else {
            set.Add(id)
        }
    }
    fmt.Println("number of id:", set.Cardinality())
}
