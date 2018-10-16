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

package types_test

import (
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestInt8rangeTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "Int8range", []interface{}{
		&types.Int8range{LowerType: types.Empty, UpperType: types.Empty, Status: types.Present},
		&types.Int8range{Lower: types.Int8{Int: 1, Status: types.Present}, Upper: types.Int8{Int: 10, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		&types.Int8range{Lower: types.Int8{Int: -42, Status: types.Present}, Upper: types.Int8{Int: -5, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		&types.Int8range{Lower: types.Int8{Int: 1, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Unbounded, Status: types.Present},
		&types.Int8range{Upper: types.Int8{Int: 1, Status: types.Present}, LowerType: types.Unbounded, UpperType: types.Exclusive, Status: types.Present},
		&types.Int8range{Status: types.Null},
	})
}

func TestInt8rangeNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select Int8range(1, 10, '(]')",
			Value: types.Int8range{Lower: types.Int8{Int: 2, Status: types.Present}, Upper: types.Int8{Int: 11, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		},
	})
}
