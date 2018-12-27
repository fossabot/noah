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

package cluster

import (
	"github.com/readystock/noah/Prototype/datums"
	"math/rand"
	"time"
)

var (
	Nodes = map[int]datums.Node{
		1: {
			NodeID:           1,
			Name:             "Node A",
			ConnectionString: "postgresql://postgres:Spring!2016@localhost:5432/ready_one?sslmode=disable",
		},
		2: {
			NodeID:           2,
			Name:             "Node B",
			ConnectionString: "postgresql://postgres:Spring!2016@localhost:5432/ready_two?sslmode=disable",
		},
		3: {
			NodeID:           3,
			Name:             "Node C",
			ConnectionString: "postgresql://postgres:Spring!2016@localhost:5432/ready_three?sslmode=disable",
		},
		4: {
			NodeID:           4,
			Name:             "Node D",
			ConnectionString: "postgresql://postgres:Spring!2016@localhost:5432/ready_four?sslmode=disable",
		},
	}

	Tables = map[string]datums.Table{
		"products": {
			TableName: "products",
			IsGlobal:  false,
		},
		"users": {
			TableName: "users",
			IsGlobal:  true,
		},
		"accounts": {
			TableName:     "accounts",
			IsGlobal:      true,
			IsTenantTable: true,
		},
	}

	Accounts = map[int]datums.Account{
		1: {
			AccountID:   1,
			AccountName: "Test Account 1",
			NodeIDs:     []int{1, 2},
		},
		2: {
			AccountID:   2,
			AccountName: "Elliot's Account",
			NodeIDs:     []int{3, 4},
		},
	}
)

func GetRandomNode() datums.Node {
	i := rand.Intn(len(Nodes))
	var v datums.Node
	for _, v = range Nodes {
		if i == 0 {
			break
		}
		i--
	}
	return v
}

func GetNodeForAccount(account_id int) datums.Node {
	rand.Seed(time.Now().Unix())
	return Nodes[Accounts[account_id].NodeIDs[rand.Intn(len(Accounts[account_id].NodeIDs))]]

}
