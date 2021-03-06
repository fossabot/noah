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

package snowflake

import (
	"errors"
	"sync"
	"time"
)

// These constants are the bit lengths of snowflake ID parts.
const (
	bitLenTime        = 39                               // bit length of time
	bitLenSequence    = 8                                // bit length of sequence number
	bitLenMachineID   = 63 - bitLenTime - bitLenSequence // bit length of machine id
	snowflakeTimeUnit = 1e7                              // nanosecond, i.e. 10 msec
)

var (
	// This is the epoch that will be used for generating unique distributed IDs.
	startTime = time.Date(2018, time.January, 9, 0, 0, 0, 0, time.UTC).UnixNano() / snowflakeTimeUnit
)

type Snowflake struct {
	mutex       *sync.Mutex
	elapsedTime int64
	sequence    uint16
	machineId   uint16
}

func NewSnowflake(machineId uint16) *Snowflake {
	sf := Snowflake{
		mutex:    new(sync.Mutex),
		sequence: uint16(1<<bitLenSequence - 1),
	}
	return &sf
}

func (sf *Snowflake) NextID() (uint64, error) {
	const maskSequence = uint16(1<<bitLenSequence - 1)
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := currentElapsedTime()
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else {
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			time.Sleep(sleepTime(sf.elapsedTime - current))
		}
	}
	return sf.toID()
}

func currentElapsedTime() int64 {
	return toSnowflakeTime(time.Now()) - startTime
}

func toSnowflakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / snowflakeTimeUnit
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond - time.Duration(time.Now().UTC().UnixNano()%snowflakeTimeUnit)*time.Nanosecond
}

func (sf *Snowflake) toID() (uint64, error) {
	if sf.elapsedTime >= 1<<bitLenTime {
		return 0, errors.New("over the time limit")
	}

	return uint64(sf.elapsedTime)<<(bitLenSequence+bitLenMachineID) |
		uint64(sf.sequence)<<bitLenMachineID |
		uint64(sf.machineId), nil
}
