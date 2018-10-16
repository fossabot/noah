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
 */

package dist_id_test

import (
	"fmt"
	"sync"
	"testing"
)

const (
	NumberOfNodes = 4
	RangeSize     = 1000
	PreRetrieve   = 60
	GenerateLoops = 4000
)

type RangeMessage struct {
	Start  uint64
	End    uint64
	Offset uint64
	Count  uint64
}

type CoordinatorNode struct {
	sequenceSync *sync.Mutex
	CurrentID    uint64
	CurrentRange *RangeMessage
	NextRange    *RangeMessage
	Master       *MasterCoordinator
}

type MasterCoordinator struct {
	CoordinatorNode
	rangeSync     *sync.Mutex
	CurrentValue  uint64
	LastNodeIndex uint64
	MaxNodeIndex  uint64
	Nodes         []*CoordinatorNode
}

func (m *MasterCoordinator) GetNextRange() RangeMessage {
	m.rangeSync.Lock()
	defer m.rangeSync.Unlock()
	if m.LastNodeIndex > m.MaxNodeIndex {
		m.CurrentValue += RangeSize
		m.LastNodeIndex = 0
		m.MaxNodeIndex = uint64(len(m.Nodes) - 1)
	}

	index := m.LastNodeIndex
	m.LastNodeIndex++

	return RangeMessage{
		Start:  m.CurrentValue,
		End:    m.CurrentValue + RangeSize,
		Offset: index,
		Count:  m.MaxNodeIndex,
	}
}

func (m *MasterCoordinator) AddNode(node *CoordinatorNode) {
	m.rangeSync.Lock()
	defer m.rangeSync.Unlock()
	node.Master = m
	m.Nodes = append(m.Nodes, node)
	if m.LastNodeIndex == 0 {
		m.MaxNodeIndex++
	}
}

func (c *CoordinatorNode) NextID() uint64 {
	c.sequenceSync.Lock()
	defer c.sequenceSync.Unlock()
	if c.CurrentRange == nil {
		cRange := c.Master.GetNextRange()
		c.CurrentRange = &cRange
	}
NewID:
	nextId := c.CurrentRange.Start + c.CurrentRange.Offset + (c.CurrentRange.Count * c.CurrentID) - (c.CurrentRange.Count - 1)
	if nextId > c.CurrentRange.End {
		// We have reached the end of our range
		if c.NextRange != nil && c.NextRange.Start >= c.CurrentRange.End {
			c.CurrentRange = c.NextRange
			c.NextRange = nil
			c.CurrentID = 1
			goto NewID
		} else {
			nRange := c.Master.GetNextRange()
			c.NextRange = &nRange
		}
	}

	if c.NextRange == nil && float64(c.CurrentID*c.CurrentRange.Count)/float64(c.CurrentRange.End-c.CurrentRange.Start) > (float64(PreRetrieve) / 100) {
		nRange := c.Master.GetNextRange()
		c.NextRange = &nRange
	}
	c.CurrentID++
	return nextId
}

func Test_GenerateDistributedIDs(t *testing.T) {
	master := MasterCoordinator{
		rangeSync:     new(sync.Mutex),
		CurrentValue:  0,
		LastNodeIndex: 0,
		MaxNodeIndex:  0,
	}
	for n := 0; n < NumberOfNodes; n++ {
		node := CoordinatorNode{
			sequenceSync: new(sync.Mutex),
			CurrentID:    1,
		}
		master.AddNode(&node)
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(master.Nodes))
	for n := 0; n < len(master.Nodes); n++ {
		go func(node *CoordinatorNode, nodeId int) {
			for i := 0; i < GenerateLoops; i++ {
				fmt.Printf("Node ID [%d] New ID: [%d]\n", nodeId, node.NextID())
			}
			wg.Done()
		}(master.Nodes[n], n)
	}
	wg.Wait()
}
