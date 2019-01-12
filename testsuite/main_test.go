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
    "github.com/jackc/pgx"
    "github.com/readystock/golog"
    "github.com/readystock/noah/cmd"
    "github.com/readystock/noah/db/system"
    "github.com/readystock/noah/testutils"
    "net"
    "os"
    "testing"
    "time"
)

var (
    SystemCtx *system.SContext
    Nodes     = []system.NNode{
        {
            Address:   "127.0.0.1",
            Port:      5432,
            Database:  "ready_one",
            User:      "postgres",
            Password:  "Spring!2016",
            ReplicaOf: 0,
            Region:    "",
            Zone:      "",
        },
    }

    Connection *pgx.Conn
)

func recoverName() {
    if r := recover(); r!= nil {
        golog.Fatalf("recovered from %v", r)
    }
}

func TestMain(m *testing.M) {
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

    for _, node := range Nodes {
        if newNode, err := SystemCtx.Nodes.AddNode(node); err != nil {
            panic(err)
        } else {
            SystemCtx.Nodes.SetNodeLive(newNode.NodeId, true)
        }
    }

    golog.Infof("finished adding %d test node(s) to cluster", len(Nodes))


    golog.Warnf("connecting to noah")
    conn := GetConnection()
    golog.Warnf("connection to noah successful")

    Connection = conn

    defer func() {
        golog.Debugf("closing noah connection")
        Connection.Close()
    }()

    func () {
        defer recoverName()
        retCode := m.Run()
        os.Exit(retCode)
    }()
    os.Exit(0)
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
