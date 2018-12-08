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
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 */

package timeofday

import (
	"math/rand"
	"time"

	"fmt"

	"strings"

	"github.com/readystock/noah/db/util/duration"
	"github.com/readystock/noah/db/util/timeutil"
)

// TimeOfDay represents a time of day (no date), stored as microseconds since
// midnight.
type TimeOfDay int64

const (
	// Min is the minimum TimeOfDay value (midnight).
	Min = TimeOfDay(0)
	// Max is the maximum TimeOfDay value (1 microsecond before midnight).
	Max = TimeOfDay(microsecondsPerDay - 1)

	microsecondsPerSecond = 1e6
	microsecondsPerMinute = 60 * microsecondsPerSecond
	microsecondsPerHour   = 60 * microsecondsPerMinute
	microsecondsPerDay    = 24 * microsecondsPerHour
	nanosPerMicro         = 1000
	secondsPerDay         = 24 * 60 * 60
)

// New creates a TimeOfDay representing the specified time.
func New(hour, min, sec, micro int) TimeOfDay {
	hours := time.Duration(hour) * time.Hour
	minutes := time.Duration(min) * time.Minute
	seconds := time.Duration(sec) * time.Second
	micros := time.Duration(micro) * time.Microsecond
	return FromInt(int64((hours + minutes + seconds + micros) / time.Microsecond))
}

func (t TimeOfDay) String() string {
	micros := t.Microsecond()
	if micros > 0 {
		s := fmt.Sprintf("%02d:%02d:%02d.%06d", t.Hour(), t.Minute(), t.Second(), micros)
		return strings.TrimRight(s, "0")
	}
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

// FromInt constructs a TimeOfDay from an int64, representing microseconds since
// midnight. Inputs outside the range [0, microsecondsPerDay) are modded as
// appropriate.
func FromInt(i int64) TimeOfDay {
	return TimeOfDay(positiveMod(i, microsecondsPerDay))
}

// positive_mod returns x mod y in the range [0, y). (Go's modulo operator
// preserves sign.)
func positiveMod(x, y int64) int64 {
	if x < 0 {
		return x%y + y
	}
	return x % y
}

// FromTime constructs a TimeOfDay from a time.Time, ignoring the date and time zone.
func FromTime(t time.Time) TimeOfDay {
	// Adjust for timezone offset so it won't affect the time. This is necessary
	// at times, like when casting from a TIMESTAMP WITH TIME ZONE.
	_, offset := t.Zone()
	unixSeconds := t.Unix() + int64(offset)

	nanos := (unixSeconds%secondsPerDay)*int64(time.Second) + int64(t.Nanosecond())
	return FromInt(nanos / nanosPerMicro)
}

// ToTime converts a TimeOfDay to a time.Time, using the Unix epoch as the date.
func (t TimeOfDay) ToTime() time.Time {
	return timeutil.Unix(0, int64(t)*nanosPerMicro)
}

// Random generates a random TimeOfDay.
func Random(rng *rand.Rand) TimeOfDay {
	return TimeOfDay(rng.Int63n(microsecondsPerDay))
}

// Add adds a Duration to a TimeOfDay, wrapping into the next day if necessary.
func (t TimeOfDay) Add(d duration.Duration) TimeOfDay {
	return FromInt(int64(t) + d.Nanos/nanosPerMicro)
}

// Difference returns the interval between t1 and t2, which may be negative.
func Difference(t1 TimeOfDay, t2 TimeOfDay) duration.Duration {
	return duration.Duration{Nanos: int64(t1-t2) * nanosPerMicro}
}

// Hour returns the hour specified by t, in the range [0, 23].
func (t TimeOfDay) Hour() int {
	return int(int64(t)%microsecondsPerDay) / microsecondsPerHour
}

// Minute returns the minute offset within the hour specified by t, in the
// range [0, 59].
func (t TimeOfDay) Minute() int {
	return int(int64(t)%microsecondsPerHour) / microsecondsPerMinute
}

// Second returns the second offset within the minute specified by t, in the
// range [0, 59].
func (t TimeOfDay) Second() int {
	return int(int64(t)%microsecondsPerMinute) / microsecondsPerSecond
}

// Microsecond returns the microsecond offset within the second specified by t,
// in the range [0, 999999].
func (t TimeOfDay) Microsecond() int {
	return int(int64(t) % microsecondsPerSecond)
}
