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

package cluster

import (
    "github.com/Ready-Stock/noah/Prototype/datums"
		"math/rand"
		"time"
)

var (
	Nodes = map[int]datums.Node {
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

	Tables = map[string]datums.Table {
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

	Accounts = map[int]datums.Account {
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