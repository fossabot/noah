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

package dist_id_test

import (
	"github.com/readystock/golog"
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

	if c.NextRange == nil && float64(c.CurrentID*c.CurrentRange.Count)/float64(c.CurrentRange.End-c.CurrentRange.Start) > (float64(PreRetrieve)/100) {
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
				golog.Verbosef("Node ID [%d] New ID: [%d]\n", nodeId, node.NextID())
			}
			wg.Done()
		}(master.Nodes[n], n)
	}
	wg.Wait()
}
