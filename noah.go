package main

import (
	"flag"
	"fmt"
	"github.com/Ready-Stock/Noah/api"
	"github.com/Ready-Stock/Noah/db/coordinator"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/badger"
	"github.com/kataras/golog"
)

type ServiceType string

const (
	Coordinator ServiceType = "coordinator"
	Tablet      ServiceType = "tablet"
)

var (
	NodeIDSequencePath = []byte("/sequences/node_id_sequence")
)

var (
	RunType       = *flag.String("type", "tablet", "Type of handler to run, defaults to `coordinator`. Valid values are: `worker` and `coordinator`.")
	HttpPort      = *flag.Int("http-port", 8080, "Listen port for Noah's HTTP REST interface.")
	PostgresPort  = *flag.Int("psql-port", 5433, "Listen port for Noah's PostgreSQL client connectivity.")
	DataDirectory = *flag.String("data-dir", "data", "Directory for Noah's embedded database.")
	WalDirectory  = *flag.String("wal-dir", "wal", "Directory for Noah's write ahead log.")
	LogLevel      = *flag.String("log-level", "debug", "Log verbosity for message written to the console window.")
)

func main() {
	flag.Parse()
	golog.SetLevel(LogLevel)

	switch ServiceType(RunType) {
	case Coordinator:
		SystemContext := system.SContext{
			Flags: system.SFlags{
				HTTPPort:      HttpPort,
				PostgresPort:  PostgresPort,
				DataDirectory: DataDirectory,
				WalDirectory:  WalDirectory,
			},
		}

		opts := badger.DefaultOptions

		opts.Dir = SystemContext.Flags.DataDirectory
		opts.ValueDir = SystemContext.Flags.DataDirectory
		badger_data, err := badger.Open(opts)
		if err != nil {
			panic(err)
		}
		node_seq, err := badger_data.GetSequence(NodeIDSequencePath, 10)
		if err != nil {
			panic(err)
		}

		opts.Dir = SystemContext.Flags.WalDirectory
		opts.ValueDir = SystemContext.Flags.WalDirectory
		badger_wal, err := badger.Open(opts)
		if err != nil {
			panic(err)
		}

		SystemContext.Badger = badger_data
		SystemContext.WalBadger = badger_wal

		defer badger_data.Close()
		defer badger_wal.Close()

		SystemContext.NodeIDs = node_seq
		fmt.Println("Starting admin application with port:", SystemContext.Flags.HTTPPort)
		fmt.Println("Listening for connections on:", SystemContext.Flags.PostgresPort)
		go coordinator.Start(&SystemContext)
		api.StartApp(&SystemContext)
	case Tablet:
		golog.Info("Starting tablet...")

	}
}
