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
 */

package system

import (
	"encoding/json"
	"fmt"
	"github.com/Ready-Stock/badger"
)

const (
	AccountsPath			   = "/accounts/"
	AccountNodesPath		   = "/account_nodes/%d/" // `/account_nodes/0000/0000`
)

type NAccount struct {
	AccountID uint64
}

func (ctx *BaseContext) GetAccounts() (a []NAccount, e error) {
	a = make([]NAccount, 0)
	e = ctx.Badger.View(func(txn *badger.Txn) error {
		accountsIterator := txn.NewIterator(badger.DefaultIteratorOptions)
		defer accountsIterator.Close()
		accountsPrefix := []byte(AccountsPath)
		for accountsIterator.Seek(accountsPrefix); accountsIterator.ValidForPrefix(accountsPrefix); accountsIterator.Next() {
			item := accountsIterator.Item()
			account := NAccount{}
			v, err := item.Value()
			if err != nil {
				return err
			}
			if err := json.Unmarshal(v, &account); err != nil {
				return err
			}
			a = append(a, account)
		}
		return nil
	})
	return a, e
}

func (ctx *BaseContext) GetNodesForAccount(accountId uint64) (n []NNode, e error) {
	e = ctx.Badger.View(func(txn *badger.Txn) error {
		n, e = ctx.getNodesForAccountEx(txn, accountId)
		return e
	})
	return n, e
}

func (ctx *BaseContext) getNodesForAccountEx(txn *badger.Txn, accountId uint64) (n []NNode, e error) {
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

