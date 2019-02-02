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

package health

import (
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/system"
	"time"
)

func StartHealthChecker(sctx *system.SContext) error {
	isLeader := sctx.IsLeader()

	for {
		time.Sleep(500 * time.Millisecond)
		// If this node gets promoted to leader
		newLeader := sctx.IsLeader()
		if !isLeader && newLeader {
			golog.Info("promoted to leader, beginning health checks")
			isLeader = newLeader
		}

		if !isLeader {
			continue
		}

		nodes, err := sctx.Nodes.GetNodes()
		if err != nil {
			golog.Errorf("could not retrieve nodes; %s", err.Error())
		}

		setNodeLive := func(nodeId uint64, alive bool) {
			if err := sctx.Nodes.SetNodeLive(nodeId, alive); err != nil {
				golog.Errorf("could not update live state of node [%d]; %s", nodeId, err.Error())
			}
		}

		for _, node := range nodes {
			func(node system.NNode) {

				conn, err := sctx.Pool.Acquire(node.NodeId)
				if err != nil {
					if node.IsAlive {
						golog.Warnf("could not connect to node [%d]; %s", node.NodeId, err.Error())
						setNodeLive(node.NodeId, false)
					}
					return
				}
				defer sctx.Pool.Release(conn)

				rows, err := conn.Query("SELECT 1;")
				if err != nil {
					if node.IsAlive {
						golog.Warnf("could not query node [%d]; %s", node.NodeId, err.Error())
						setNodeLive(node.NodeId, false)
					}
					return
				}
				defer rows.Close()

				result := rows.Next()
				if err := rows.Err(); err != nil {
					if node.IsAlive {
						golog.Warnf("could not query node [%d]; %s", node.NodeId, err.Error())
						setNodeLive(node.NodeId, false)
					}
				} else if result {
					if !node.IsAlive {
						golog.Infof("node [%d] is now alive", node.NodeId)
						setNodeLive(node.NodeId, true)
					}
				}
			}(node)
		}
	}
}
