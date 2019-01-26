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

package system

import (
	"github.com/ahmetb/go-linq"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/go-errors"
)

type SAccounts baseContext

func (ctx *SAccounts) GetAccount(accountId uint64) (account *NAccount, err error) {
	accountBytes, err := ctx.db.Get(getAccountPath(accountId))
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(accountBytes, account)
	return account, err
}

func (ctx *SAccounts) GetAccounts() (accounts []NAccount, err error) {
	accountsBytes, err := ctx.db.GetPrefix(getAccountsPath())
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
	sNode := SNode(*ctx)
	sSettings := SSettings(*ctx)
	sSequence := SSequence(*ctx)
	nodes, err := (&sNode).GetNodes()
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
	replicationFactor, err := (&sSettings).GetSettingInt64(QueryReplicationFactor)
	if err != nil {
		return nil, nil, err
	}

	if int64(len(liveNodes)) < *replicationFactor { // If there are not enough nodes to adequately replicate data.
		return nil, nil, errors.New("could not create account, replication factor is greater than the number of nodes available in cluster")
	}

	accountId, err := (&sSequence).NewAccountID()
	if err != nil {
		return nil, nil, err
	}

	accountNodes := make([]NNode, *replicationFactor)
	if int64(len(liveNodes)) == *replicationFactor {
		accountNodes = liveNodes
		for i := 0; i < len(liveNodes); i++ {
			err = ctx.db.Set(getAccountsNodesAccountNodePath(*accountId, accountNodes[i].NodeId), []byte{})
			if err != nil {
				return nil, nil, err
			}
		}
	} else {
		for i := uint64(0); i < uint64(*replicationFactor); i++ {
			accountNodes[i] = liveNodes[(*accountId+(i*uint64(*replicationFactor)))%uint64(len(nodes))] // This math takes the account id and distributes it in the cluster to pick a node
			// TODO (elliotcourant) if the replication factor is > 1 and at least 1 of these sets fail then it could mess up the records of what nodes host what account
			err = ctx.db.Set(getAccountsNodesAccountNodePath(*accountId, accountNodes[i].NodeId), []byte{})
			if err != nil {
				return nil, nil, err
			}
		}
	}

	account := &NAccount{
		AccountId: *accountId,
	}
	b, err := proto.Marshal(account)
	if err != nil {
		return nil, nil, err
	}

	err = ctx.db.Set(getAccountsNodesAccountPath(*accountId), b)
	if err != nil {
		return nil, nil, err
	}
	return account, accountNodes, nil
}

func (ctx *SAccounts) GetNodesForAccounts(accountIds ...uint64) (nodes []NNode, err error) {
	nodes = make([]NNode, 0)
	for _, accountId := range accountIds {
		if accountNodes, err := ctx.GetNodesForAccount(accountId); err != nil {
			return nil, err
		} else {
			nodes = append(nodes, accountNodes...)
		}
	}
	return nodes, nil
}

func (ctx *SAccounts) GetNodesForAccount(accountId uint64) (nodes []NNode, err error) {
	nodeBytes, err := ctx.db.GetPrefix(getAccountsNodesAccountPath(accountId))
	if err != nil {
		return nil, err
	}
	nodes = make([]NNode, len(nodeBytes))
	for i, kv := range nodeBytes {
		node := NNode{}
		err := proto.Unmarshal(kv.Value, &node)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}
