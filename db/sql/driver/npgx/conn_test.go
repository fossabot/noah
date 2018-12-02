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
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 */

package npgx_test

import (
	"github.com/Ready-Stock/noah/db/sql/driver"
	"github.com/Ready-Stock/noah/db/sql/driver/npgx"
	"testing"
)

var (
	ConnectionConfig = driver.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "Spring!2016",
		Database: "postgres",
	}
	ConnectionConfig2 = driver.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "Spring!2016",
		Database: "ready_three",
	}

	BadConnectionConfig = driver.ConnConfig{
		Host:     "localhost",
		Port:     123,
		User:     "postgres",
		Password: "Spring!2016",
		Database: "postgres",
	}
)

func Test_Connect(t *testing.T) {
	d, err := npgx.Connect(ConnectionConfig)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	r, err := d.Query("SELECT 1;")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	for r.Next() {
		fields, err := r.Values()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if len(fields) == 1 {
			if fields[0].(int32) != 1 {
				t.Error("Unexpected result!")
				t.Fail()
			}
		}
	}
	err = d.Close()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_CITEXT(t *testing.T) {
	d, err := npgx.Connect(ConnectionConfig)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	r, err := d.Query("SELECT 123::CITEXT;")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	for r.Next() {
		fields, err := r.Values()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if len(fields) == 1 {
			if fields[0].(string) != "123" {
				t.Error("Unexpected result!")
				t.Fail()
			}
		}
	}
	err = d.Close()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}


func Test_BadConnection(t *testing.T) {
	_, err := npgx.Connect(BadConnectionConfig)
	if err == nil {
		t.Fail()
	}
}
