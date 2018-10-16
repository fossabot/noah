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

package system

import (
	"encoding/json"
	"fmt"
	"github.com/Ready-Stock/badger"
	"github.com/ahmetb/go-linq"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/go-errors"
)

const (
	AccountNodesPath = "/account_nodes/%d/" // `/account_nodes/0000/0000`
)

type SAccounts baseContext

func (ctx *SAccounts) GetAccounts() (accounts []NAccount, err error) {
	accountsBytes, err := ctx.db.GetPrefix([]byte(accountsPath))
	if err != nil {
		return nil, err
	}
	accounts = make([]NAccount, len(accountsBytes))
	for i, kv := range accountsBytes {
		account := NAccount{}
		err := proto.Unmarshal(kv.Value, &account)
		if err != nil {
			return nil, err
		}
		accounts[i] = account
	}
	return accounts, nil
}

func (ctx *SAccounts) CreateAccount() (*NAccount, []NNode, error) {
	nodes, err := SNode(*ctx).GetNodes()
	if err != nil {
		return nil, nil, err
	}
	liveNodes := make([]NNode, 0)
	linq.From(nodes).WhereT(func(node NNode) bool {
		return node.IsAlive && node.ReplicaOf == 0
	}).ToSlice(&liveNodes)
	nonReplicas := linq.From(nodes).CountWithT(func(node NNode) bool {
		return node.ReplicaOf == 0
	})
	if len(liveNodes) < nonReplicas { // If there are any database nodes that are currently offline and accept writes, then reject the new account
		return nil, nil, errors.New("could not create account at this time, insufficient database nodes are available")
	}
	replicationFactor, err := SSettings(*ctx).GetSettingInt64(QueryReplicationFactor)
	if err != nil {
		return nil, nil, err
	}
	if int64(len(liveNodes)) < *replicationFactor { // If there are not enough nodes to adequately replicate data.
		return nil, nil, errors.New("could not create account, replication factor is greater than the number of nodes available in cluster")
	}

	accountId, err := SSequence(*ctx).NewAccountID()
	if err != nil {
		return nil, nil, err
	}
	accountNodes := make([]NNode, *replicationFactor)
	for i := uint64(0); i < uint64(*replicationFactor); i++ {
		accountNodes[i] = liveNodes[(*accountId + (i * uint64(*replicationFactor))) % uint64(len(nodes))]
		// TODO (elliotcourant) if the replication factor is > 1 and at least 1 of these sets fail then it could mess up the records of what nodes host what account
		err = ctx.db.Set([]byte(fmt.Sprintf("%s%d/%d", accountsNodesPath, accountId, accountNodes[i].NodeId)), []byte{})
		if err != nil {
			return nil, nil, err
		}
	}

	account := &NAccount{
		AccountId:*accountId,
	}
	b, err := proto.Marshal(account)
	if err != nil {
		return nil, nil, err
	}

	err = ctx.db.Set([]byte(fmt.Sprintf("%s%d", accountsPath, accountId)), b)
	if err != nil {
		return nil, nil, err
	}
	return account, accountNodes, nil
}

func (ctx *SAccounts) GetNodesForAccount(accountId uint64) (n []NNode, e error) {
	e = ctx.Badger.View(func(txn *badger.Txn) error {
		n, e = ctx.getNodesForAccountEx(txn, accountId)
		return e
	})
	return n, e
}

func (ctx *SAccounts) getNodesForAccountEx(txn *badger.Txn, accountId uint64) (n []NNode, e error) {
	nodesIterator := txn.NewIterator(badger.DefaultIteratorOptions)
	defer nodesIterator.Close()
	n = make([]NNode, 0)
	accountNodesPrefix := []byte(fmt.Sprintf(AccountNodesPath, accountId))
	for nodesIterator.Seek(accountNodesPrefix); nodesIterator.ValidForPrefix(accountNodesPrefix); nodesIterator.Next() {
		nodeItem := nodesIterator.Item()
		nodeValue, err := nodeItem.Value()
		if err != nil {
			return nil, err
		}
		node := NNode{}
		if err := json.Unmarshal(nodeValue, &node); err != nil {
			return nil, err
		}
		n = append(n, node)
	}
	return n, nil
}
