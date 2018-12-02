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

package sessiondata

import (
    "github.com/Ready-Stock/noah/db/sql/pgwire/pgerror"
    "github.com/Ready-Stock/noah/db/util/syncutil"
)

// SequenceState stores session-scoped state used by sequence builtins.
//
// All public methods of SequenceState are thread-safe, as the structure is
// meant to be shared by statements executing in parallel on a session.
type SequenceState struct {
	mu struct {
		syncutil.Mutex
		// latestValues stores the last value obtained by nextval() in this session
		// by descriptor id.
		latestValues map[uint32]int64

		// lastSequenceIncremented records the descriptor id of the last sequence
		// nextval() was called on in this session.
		lastSequenceIncremented uint32
	}
}

// NewSequenceState creates a SequenceState.
func NewSequenceState() *SequenceState {
	ss := SequenceState{}
	ss.mu.latestValues = make(map[uint32]int64)
	return &ss
}

// copy performs a deep copy of SequenceState.
func (ss *SequenceState) copy() *SequenceState {
	cp := NewSequenceState()
	ss.mu.Lock()
	defer ss.mu.Unlock()
	for k, v := range ss.mu.latestValues {
		cp.mu.latestValues[k] = v
	}
	cp.mu.lastSequenceIncremented = ss.mu.lastSequenceIncremented
	return ss
}

// NextVal ever called returns true if a sequence has ever been incremented on
// this session.
func (ss *SequenceState) nextValEverCalledLocked() bool {
	return len(ss.mu.latestValues) > 0
}

// RecordValue records the latest manipulation of a sequence done by a session.
func (ss *SequenceState) RecordValue(seqID uint32, val int64) {
	ss.mu.Lock()
	ss.mu.lastSequenceIncremented = seqID
	ss.mu.latestValues[seqID] = val
	ss.mu.Unlock()
}

// SetLastSequenceIncremented sets the id of the last incremented sequence.
// Usually this id is set through RecordValue().
func (ss *SequenceState) SetLastSequenceIncremented(seqID uint32) {
	ss.mu.Lock()
	ss.mu.lastSequenceIncremented = seqID
	ss.mu.Unlock()
}

// GetLastValue returns the value most recently obtained by
// nextval() for the last sequence for which RecordLatestVal() was called.
func (ss *SequenceState) GetLastValue() (int64, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if !ss.nextValEverCalledLocked() {
		return 0, pgerror.NewError(
			pgerror.CodeObjectNotInPrerequisiteStateError, "lastval is not yet defined in this session")
	}

	return ss.mu.latestValues[ss.mu.lastSequenceIncremented], nil
}

// GetLastValueByID returns the value most recently obtained by nextval() for
// the given sequence in this session.
// The bool retval is false if RecordLatestVal() was never called on the
// requested sequence.
func (ss *SequenceState) GetLastValueByID(seqID uint32) (int64, bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	val, ok := ss.mu.latestValues[seqID]
	return val, ok
}

// Export returns a copy of the SequenceState's state - the latestValues and
// lastSequenceIncremented.
// lastSequenceIncremented is only defined if latestValues is non-empty.
func (ss *SequenceState) Export() (map[uint32]int64, uint32) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	res := make(map[uint32]int64, len(ss.mu.latestValues))
	for k, v := range ss.mu.latestValues {
		res[k] = v
	}
	return res, ss.mu.lastSequenceIncremented
}
