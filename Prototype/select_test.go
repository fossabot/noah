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
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

package Prototype

import (
	"fmt"
	"testing"
)

var (
	SelectQueries = []string{
		"SELECT 1",
		"SELECT 1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id=1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id = 1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id = '1' OR product_id = '2';",
		"SELECT products.product_id FROM products WHERE (account_id = '1') AND (product_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id FROM products WHERE (products.account_id = '1') AND (products.product_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id WHERE (account_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id INNER JOIN temp ON temp.id=products.product_id WHERE (account_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id INNER JOIN temp ON temp.id=products.product_id WHERE (product.id = 1 or product.test = true) and (products.account_id = '1') LIMIT 10 OFFSET 0;",
	}
	FunctionQueries = []string{
		"SELECT do_test('test');",
		"SELECT CURRENT_TIMESTAMP;",
	}
)

func Test_Select1(t *testing.T) {
	context := Start()
	if rows, err := InjestQuery(&context, "SELECT 1;"); err != nil {
		t.Error(err)
	} else {
		for rows.Next() {
			var col = ""
			rows.Scan(&col)
			fmt.Printf("Row (%s)", col)
		}
	}
}

func Test_Selects(t *testing.T) {
	for _, q := range SelectQueries {
		context := Start()
		if _, err := InjestQuery(&context, q); err != nil {
			t.Error(err)
		}
	}
}

func Test_SelectFunctions(t *testing.T) {
	for _, q := range FunctionQueries {
		context := Start()
		if _, err := InjestQuery(&context, q); err != nil {
			t.Error(err)
		}
	}
}

func Benchmark_Select(b *testing.B) {
	context := Start()
	if _, err := InjestQuery(&context, "SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id INNER JOIN temp ON temp.id=products.product_id WHERE (product.id = 1 or product.test = true) and (products.account_id = '1') LIMIT 10 OFFSET 0;"); err != nil {
		b.Error(err)
	}
}
