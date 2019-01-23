/*
 * Copyright (c) 2019 Ready Stock
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

package testsuite

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/readystock/golog"
	"github.com/readystock/noah/cmd"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/noah/testutils"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
	"time"
)

var (
	postgresConfig = pgx.ConnConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "postgres",
		User:     "postgres",
		Password: "Spring!2016",
		LogLevel: 6,
		Logger:   NewLogger(),
	}
	SystemCtx     *system.SContext
	NumberOfNodes = 1
	Connection    *pgx.Conn
)

func recoverName() {
	if r := recover(); r != nil {
		golog.Fatalf("recovered from %v", r)
	}
}

func TestMain(m *testing.M) {
	golog.Infof("Testing with %d node(s)", NumberOfNodes)
	nodes := make([]system.NNode, NumberOfNodes)
	for i := 0; i < NumberOfNodes; i++ {
		nodes[i] = system.NNode{
			Address:   "127.0.0.1",
			Port:      5432,
			Database:  fmt.Sprintf("ready_test_%d", i+1),
			User:      "postgres",
			Password:  "Spring!2016",
			ReplicaOf: 0,
			Region:    "",
			Zone:      "",
		}
	}

	golog.SetLevel("trace")
	retCode := func() int {
		golog.Infof("SETTING UP TEST DATABASE")

		postgres, err := pgx.Connect(postgresConfig)
		if err != nil {
			panic(err)
		}

		for _, node := range nodes {
			if _, err := postgres.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", node.Database)); err != nil {
				panic(err)
			}

			if _, err := postgres.Exec(fmt.Sprintf("CREATE DATABASE %s", node.Database)); err != nil {
				panic(err)
			}
		}

		defer func() {
			golog.Infof("TEARING DOWN TEST DATABASES")
			for _, node := range nodes {
				if _, err := postgres.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s';", node.Database)); err != nil {
					panic(err)
				}

				if _, err := postgres.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", node.Database)); err != nil {
					panic(err)
				}
			}
		}()

		golog.Info("RUNNING SQL TEST")
		tempFolder := testutils.CreateTempFolder()
		defer testutils.DeleteTempFolder(tempFolder)

		sctx := new(system.SContext)
		go func() {
			cmd.RunNode(cmd.NodeConfig{
				PGWireAddr:     "127.0.0.1:0",
				GrpcAddr:       "127.0.0.1:0",
				WebAddr:        "127.0.0.1:0",
				JoinAddr:       "",
				StoreDirectory: tempFolder,
			}, sctx)
		}()
		golog.Info("waiting for cluster to start up")

		for !sctx.IsRunning() {
			time.Sleep(5 * time.Second)
		}

		golog.Info("finished setting up test cluster")

		SystemCtx = sctx

		for _, node := range nodes {
			if _, err := SystemCtx.Nodes.AddNode(node); err != nil {
				panic(err)
			} else {
				// SystemCtx.Nodes.SetNodeLive(newNode.NodeId, true)
			}
		}

		// We want to wait just a short while while the health checker verifies the new nodes.
		time.Sleep(5 * time.Second)

		golog.Infof("finished adding %d test node(s) to cluster", len(nodes))

		golog.Warnf("connecting to noah")
		conn := GetConnection()
		golog.Warnf("connection to noah successful")

		Connection = conn

		defer func() {
			if err := Connection.Close(); err != nil {
				golog.Fatalf("error closing pgx connection, %v", err)
			}
		}()

		defer recoverName()
		return m.Run()
	}()
	os.Exit(retCode)
}

func GetConnection() *pgx.Conn {
	addr, err := net.ResolveTCPAddr("tcp", SystemCtx.PgListenAddr())
	if err != nil {
		panic(err)
	}
	golog.Infof("connecting to noah at address `%s`", SystemCtx.PgListenAddr())

	config := pgx.ConnConfig{
		Host:     addr.IP.String(),
		Port:     uint16(addr.Port),
		Database: "postgres",
		User:     "postgres",
		Password: "password",
		LogLevel: 6,
		Logger:   NewLogger(),
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		panic(err)
	}
	return conn
}

func DoQueryTest(t *testing.T, test QueryTest) [][]interface{} {
	startTime := time.Now()
	defer func() {
		golog.Tracef("FINISHED TESTING QUERY, TIME: %v", time.Since(startTime))
	}()
	golog.Tracef("SENDING QUERY `%s`", test.Query)
	result, err := Connection.Query(test.Query, test.Args...)
	if err != nil {
		golog.Error(result.Err())
		t.Error(err)
		t.FailNow()
	}

	defer result.Close()

	results := make([][]interface{}, len(test.Expected))

	index := 0
	for result.Next() {
		if result.Err() != nil {
			golog.Error(result.Err())
			t.Error(err)
			t.FailNow()
		}

		vals, err := result.Values()
		if err != nil {
			golog.Error(result.Err())
			t.Error(err)
			t.FailNow()
		}

		if test.Expected != nil {
			assert.EqualValuesf(t, test.Expected[index], vals,
				"`%s` | %v - row [%d] did not return expected value", test.Query, test.Args, index)
		}

		results = append(results, vals)
		index++
	}

	if result.Err() != nil {
		golog.Error(result.Err())
		t.Error(err)
		t.FailNow()
	}

	if test.Expected != nil {
		assert.Equal(t, len(test.Expected), index, "`%s` | %v - number of rows returned did not match expected", test.Query, test.Args)
	}

	return results
}

func DoExecTest(t *testing.T, test ExecTest) {
	result, err := Connection.Exec(test.Query, test.Args...)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, test.Expected, result.RowsAffected(),
		"`%s` | %v - did not return expected rows affected", test.Query, test.Args)
}
