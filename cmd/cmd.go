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

package cmd

import (
	"fmt"
	"github.com/readystock/golog"
	"github.com/readystock/noah/api"
	"github.com/readystock/noah/db/coordinator"
	"github.com/readystock/noah/db/health"
	"github.com/readystock/noah/db/system"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	systemContext *system.SContext
	ch            chan os.Signal
)

var (
	rootCmd = &cobra.Command{
		Use:   "noah",
		Short: "Noah is a distributed new-SQL database.",
		Long: `
Noah is a distributed new-SQL database designed for multi-tenant
applications. It is powered by a sharded PostgreSQL backend. 

To start simply run the command 'noah'. This will start a single
coordinator instance of noah. From there you will want to add a
PostgreSQL node via the HTTP or gRPC interface. Once a database
node has been added, you can issue queries to create tables and
build your database. 

To join an existing cluster, specify the address of any node in
the existing cluster with the --join flag. This will send a 
message to that node and get any needed information about the 
leader of the cluster and perform the actions necessary to join.

To connect to noah, use a standard PostgreSQL connection string
with the user 'root' and no password; use port 5433 unless you
have specified a custom address for noah to listen to for pgwire
connections. This method accepts pgwire protocol and should work
with almost all typical Postgres drivers and IDEs.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	startCmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			RunNode(NodeConfig{
				PGWireAddr:     PGWireAddr,
				GrpcAddr:       GrpcAddr,
				WebAddr:        WebAddr,
				JoinAddr:       JoinAddr,
				StoreDirectory: StoreDirectory,
				LogLevel:       LogLevel,
			}, nil)
		},
	}
)

var (
	PGWireAddr     string
	GrpcAddr       string
	WebAddr        string
	JoinAddr       string
	StoreDirectory string
	LogLevel       string
)

func init() {
	startCmd.Flags().StringVarP(&PGWireAddr, "pgwire", "p", ":5433", "address that will accept PostgreSQL connections")
	startCmd.Flags().StringVarP(&GrpcAddr, "grpc", "g", ":5434", "address that will be used for Noah's raft/gRPC interface")
	startCmd.Flags().StringVarP(&WebAddr, "web", "w", ":5435", "address that will be used for Noah's HTTP interface")
	startCmd.Flags().StringVarP(&JoinAddr, "join", "j", "", "address and gRPC port of another node in a cluster to join")
	startCmd.Flags().StringVarP(&StoreDirectory, "store", "s", "data", "directory that will be used for Noah's key value store")
	startCmd.Flags().StringVarP(&LogLevel, "log", "l", "verbose", "log output level, valid values: trace, verbose, debug, info, warn, error, fatal")
	rootCmd.AddCommand(startCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type NodeConfig struct {
	PGWireAddr     string
	GrpcAddr       string
	WebAddr        string
	JoinAddr       string
	StoreDirectory string
	LogLevel       string
}

func RunNode(config NodeConfig, sctx *system.SContext) {
	golog.SetLevel(config.LogLevel)
	golog.Infof("Starting noah...")

	s, err := system.NewSystemContext(config.StoreDirectory, config.GrpcAddr, config.JoinAddr, config.PGWireAddr)
	if err != nil {
		golog.Fatalf("A fatal error was encountered trying to initialize store; %s", err.Error())
	}
	if sctx == nil {
		sctx = s
	} else {
		*sctx = *s
	}

	golog.Infof("Coordinator ID [%d] starting...", sctx.CoordinatorID())

	systemContext = sctx
	ch = make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGSEGV)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)

	go func() {
		<-ch
		golog.Warnf("Stopping coordinator node.")
		systemContext.Close()
		os.Exit(0)
	}()

	if err := sctx.Setup.DoSetup(); err != nil {
		golog.Fatalf("failed to perform setup on node; %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := coordinator.Start(sctx); err != nil {
			golog.Error(err)
		}
	}()
	go func() {
		defer wg.Done()
		api.StartApp(sctx, config.WebAddr)
	}()
	go func() {
		defer wg.Done()
		health.StartHealthChecker(sctx)
	}()
	golog.Infof("Noah is now running.")
	sctx.FinishedStartup()
	wg.Wait()
}
