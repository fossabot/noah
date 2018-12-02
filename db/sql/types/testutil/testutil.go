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

package testutil

import (
	"context"
	"database/sql"
	"fmt"
    "github.com/Ready-Stock/noah/db/sql/driver/npgx"
	"github.com/Ready-Stock/pgx"
	"os"
	"reflect"
	"testing"

    "github.com/Ready-Stock/noah/db/sql/types"
	_ "github.com/lib/pq"
)

func MustConnectDatabaseSQL(t testing.TB, driverName string) *sql.DB {
	var sqlDriverName string
	switch driverName {
	case "github.com/lib/pq":
		sqlDriverName = "postgres"
	case "github.com/jackc/pgx/stdlib":
		sqlDriverName = "pgx"
	default:
		t.Fatalf("Unknown driver %v", driverName)
	}

	db, err := sql.Open(sqlDriverName, os.Getenv("PGX_TEST_DATABASE"))
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func MustConnectPgx(t testing.TB) *npgx.Conn {
	config, err := npgx.ParseConnectionString(os.Getenv("PGX_TEST_DATABASE"))
	if err != nil {
		t.Fatal(err)
	}

	conn, err := npgx.Connect(config)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func MustClose(t testing.TB, conn interface {
	Close() error
}) {
	err := conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

type forceTextEncoder struct {
	e types.TextEncoder
}

func (f forceTextEncoder) EncodeText(ci *types.ConnInfo, buf []byte) ([]byte, error) {
	return f.e.EncodeText(ci, buf)
}

type forceBinaryEncoder struct {
	e types.BinaryEncoder
}

func (f forceBinaryEncoder) EncodeBinary(ci *types.ConnInfo, buf []byte) ([]byte, error) {
	return f.e.EncodeBinary(ci, buf)
}

func ForceEncoder(e interface{}, formatCode int16) interface{} {
	switch formatCode {
	case pgx.TextFormatCode:
		if e, ok := e.(types.TextEncoder); ok {
			return forceTextEncoder{e: e}
		}
	case pgx.BinaryFormatCode:
		if e, ok := e.(types.BinaryEncoder); ok {
			return forceBinaryEncoder{e: e.(types.BinaryEncoder)}
		}
	}
	return nil
}

func TestSuccessfulTranscode(t testing.TB, typesName string, values []interface{}) {
	TestSuccessfulTranscodeEqFunc(t, typesName, values, func(a, b interface{}) bool {
		return reflect.DeepEqual(a, b)
	})
}

func TestSuccessfulTranscodeEqFunc(t testing.TB, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	TestPgxSuccessfulTranscodeEqFunc(t, typesName, values, eqFunc)
	TestPgxSimpleProtocolSuccessfulTranscodeEqFunc(t, typesName, values, eqFunc)
	for _, driverName := range []string{"github.com/lib/pq", "github.com/jackc/pgx/stdlib"} {
		TestDatabaseSQLSuccessfulTranscodeEqFunc(t, driverName, typesName, values, eqFunc)
	}
}

func TestPgxSuccessfulTranscodeEqFunc(t testing.TB, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectPgx(t)
	defer MustClose(t, conn)

	ps, err := conn.Prepare("test", fmt.Sprintf("select $1::%s", typesName))
	if err != nil {
		t.Fatal(err)
	}

	formats := []struct {
		name       string
		formatCode int16
	}{
		{name: "TextFormat", formatCode: pgx.TextFormatCode},
		{name: "BinaryFormat", formatCode: pgx.BinaryFormatCode},
	}

	for i, v := range values {
		for _, fc := range formats {
			ps.FieldDescriptions[0].FormatCode = fc.formatCode
			vEncoder := ForceEncoder(v, fc.formatCode)
			if vEncoder == nil {
				t.Logf("Skipping: %#v does not implement %v", v, fc.name)
				continue
			}
			// Derefence value if it is a pointer
			derefV := v
			refVal := reflect.ValueOf(v)
			if refVal.Kind() == reflect.Ptr {
				derefV = refVal.Elem().Interface()
			}

			result := reflect.New(reflect.TypeOf(derefV))
			err := conn.QueryRow("test", ForceEncoder(v, fc.formatCode)).Scan(result.Interface())
			if err != nil {
				t.Errorf("%v %d: %v", fc.name, i, err)
			}

			if !eqFunc(result.Elem().Interface(), derefV) {
				t.Errorf("%v %d: expected %v, got %v", fc.name, i, derefV, result.Elem().Interface())
			}
		}
	}
}

func TestPgxSimpleProtocolSuccessfulTranscodeEqFunc(t testing.TB, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectPgx(t)
	defer MustClose(t, conn)

	for i, v := range values {
		// Derefence value if it is a pointer
		derefV := v
		refVal := reflect.ValueOf(v)
		if refVal.Kind() == reflect.Ptr {
			derefV = refVal.Elem().Interface()
		}

		result := reflect.New(reflect.TypeOf(derefV))
		err := conn.QueryRowEx(
			context.Background(),
			fmt.Sprintf("select ($1)::%s", typesName),
			&pgx.QueryExOptions{SimpleProtocol: true},
			v,
		).Scan(result.Interface())
		if err != nil {
			t.Errorf("Simple protocol %d: %v", i, err)
		}

		if !eqFunc(result.Elem().Interface(), derefV) {
			t.Errorf("Simple protocol %d: expected %v, got %v", i, derefV, result.Elem().Interface())
		}
	}
}

func TestDatabaseSQLSuccessfulTranscodeEqFunc(t testing.TB, driverName, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectDatabaseSQL(t, driverName)
	defer MustClose(t, conn)

	ps, err := conn.Prepare(fmt.Sprintf("select $1::%s", typesName))
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range values {
		// Derefence value if it is a pointer
		derefV := v
		refVal := reflect.ValueOf(v)
		if refVal.Kind() == reflect.Ptr {
			derefV = refVal.Elem().Interface()
		}

		result := reflect.New(reflect.TypeOf(derefV))
		err := ps.QueryRow(v).Scan(result.Interface())
		if err != nil {
			t.Errorf("%v %d: %v", driverName, i, err)
		}

		if !eqFunc(result.Elem().Interface(), derefV) {
			t.Errorf("%v %d: expected %v, got %v", driverName, i, derefV, result.Elem().Interface())
		}
	}
}

type NormalizeTest struct {
	SQL   string
	Value interface{}
}

func TestSuccessfulNormalize(t testing.TB, tests []NormalizeTest) {
	TestSuccessfulNormalizeEqFunc(t, tests, func(a, b interface{}) bool {
		return reflect.DeepEqual(a, b)
	})
}

func TestSuccessfulNormalizeEqFunc(t testing.TB, tests []NormalizeTest, eqFunc func(a, b interface{}) bool) {
	TestPgxSuccessfulNormalizeEqFunc(t, tests, eqFunc)
	for _, driverName := range []string{"github.com/lib/pq", "github.com/jackc/pgx/stdlib"} {
		TestDatabaseSQLSuccessfulNormalizeEqFunc(t, driverName, tests, eqFunc)
	}
}

func TestPgxSuccessfulNormalizeEqFunc(t testing.TB, tests []NormalizeTest, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectPgx(t)
	defer MustClose(t, conn)

	formats := []struct {
		name       string
		formatCode int16
	}{
		{name: "TextFormat", formatCode: pgx.TextFormatCode},
		{name: "BinaryFormat", formatCode: pgx.BinaryFormatCode},
	}

	for i, tt := range tests {
		for _, fc := range formats {
			psName := fmt.Sprintf("test%d", i)
			ps, err := conn.Prepare(psName, tt.SQL)
			if err != nil {
				t.Fatal(err)
			}

			ps.FieldDescriptions[0].FormatCode = fc.formatCode
			if ForceEncoder(tt.Value, fc.formatCode) == nil {
				t.Logf("Skipping: %#v does not implement %v", tt.Value, fc.name)
				continue
			}
			// Derefence value if it is a pointer
			derefV := tt.Value
			refVal := reflect.ValueOf(tt.Value)
			if refVal.Kind() == reflect.Ptr {
				derefV = refVal.Elem().Interface()
			}

			result := reflect.New(reflect.TypeOf(derefV))
			err = conn.QueryRow(psName).Scan(result.Interface())
			if err != nil {
				t.Errorf("%v %d: %v", fc.name, i, err)
			}

			if !eqFunc(result.Elem().Interface(), derefV) {
				t.Errorf("%v %d: expected %v, got %v", fc.name, i, derefV, result.Elem().Interface())
			}
		}
	}
}

func TestDatabaseSQLSuccessfulNormalizeEqFunc(t testing.TB, driverName string, tests []NormalizeTest, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectDatabaseSQL(t, driverName)
	defer MustClose(t, conn)

	for i, tt := range tests {
		ps, err := conn.Prepare(tt.SQL)
		if err != nil {
			t.Errorf("%d. %v", i, err)
			continue
		}

		// Derefence value if it is a pointer
		derefV := tt.Value
		refVal := reflect.ValueOf(tt.Value)
		if refVal.Kind() == reflect.Ptr {
			derefV = refVal.Elem().Interface()
		}

		result := reflect.New(reflect.TypeOf(derefV))
		err = ps.QueryRow().Scan(result.Interface())
		if err != nil {
			t.Errorf("%v %d: %v", driverName, i, err)
		}

		if !eqFunc(result.Elem().Interface(), derefV) {
			t.Errorf("%v %d: expected %v, got %v", driverName, i, derefV, result.Elem().Interface())
		}
	}
}
