package cluster

import (
	data "github.com/Ready-Stock/Noah/Prototype/datums"
		"math/rand"
		"time"
)

var (
	Nodes = map[int]data.Node {
		1: {
			NodeID: 1,
			Name: "Node A",
			ConnectionString:"postgresql://postgres:Spring!2016@localhost:5432/ready_one?sslmode=disable",
		},
		2: {
			NodeID: 2,
			Name: "Node B",
			ConnectionString:"postgresql://postgres:Spring!2016@localhost:5432/ready_two?sslmode=disable",
		},
		3: {
			NodeID: 3,
			Name: "Node C",
			ConnectionString:"postgresql://postgres:Spring!2016@localhost:5432/ready_three?sslmode=disable",
		},
		4: {
			NodeID: 4,
			Name: "Node D",
			ConnectionString:"postgresql://postgres:Spring!2016@localhost:5432/ready_four?sslmode=disable",
		},
	}

	Tables = map[string]data.Table {
		"products": {
			TableName:"products",
			IsGlobal:false,
		},
		"users": {
			TableName:"users",
			IsGlobal:true,
		},
		"accounts": {
			TableName:"accounts",
			IsGlobal:true,
			IsTenantTable:true,
		},
	}

	Accounts = map[int]data.Account {
		1: {
			AccountID: 1,
			AccountName: "Test Account 1",
			NodeIDs: []int{ 1, 2 },
		},
		2: {
			AccountID: 2,
			AccountName: "Elliot's Account",
			NodeIDs: []int{ 3, 4 },
		},
	}
)

func GetRandomNode() data.Node {
	i := rand.Intn(len(Nodes))
	var v data.Node
	for _, v = range Nodes {
		if i == 0 {
			break
		}
		i--
	}
	return v
}

func GetNodeForAccount(account_id int) data.Node {
	rand.Seed(time.Now().Unix())
	return Nodes[Accounts[account_id].NodeIDs[rand.Intn(len(Accounts[account_id].NodeIDs))]]

}