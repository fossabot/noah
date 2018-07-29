package cluster

import "github.com/Ready-Stock/Noah/Prototype"

var (
	Nodes = map[int64]Node {
		1: Node {
			NodeID: 1,
			Name: "Node A",
		},
		2: Node {
			NodeID: 2,
			Name: "Node B",
		},
		3: Node {
			NodeID: 3,
			Name: "Node C",
		},
		4: Node {
			NodeID: 4,
			Name: "Node D",
		},
	}
)