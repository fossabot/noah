package main

import (
	"flag"
	"fmt"
	"github.com/Ready-Stock/Noah/db"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/badger"
	"github.com/kataras/golog"
)

type WorkerType string

const (
	Coordinator WorkerType = "coordinator"
	Worker      WorkerType = "worker"
)

var (
	NodeIDSequencePath = []byte("/node_id_sequence")
)

var (
	RunType       = *flag.String("type", "coordinator", "Type of handler to run, defaults to `coordinator`. Valid values are: `worker` and `coordinator`.")
	HttpPort      = *flag.Int("http-port", 8080, "Listen port for Noah's HTTP REST interface.")
	PostgresPort  = *flag.Int("psql-port", 5433, "Listen port for Noah's PostgreSQL client connectivity.")
	DataDirectory = *flag.String("data-dir", "data", "Directory for Noah's embedded database.")
	WalDirectory  = *flag.String("wal-dir", "wal", "Directory for Noah's write ahead log.")
	LogLevel = *flag.String("log-level", "debug", "Log verbosity for message written to the console window.")
)

func main() {
	flag.Parse()
	golog.SetLevel(LogLevel)

	SystemContext := system.SContext{
		Flags: system.SFlags{
			HTTPPort:      HttpPort,
			PostgresPort:  PostgresPort,
			DataDirectory: DataDirectory,
		},
	}


	opts := badger.DefaultOptions
	opts.Dir = SystemContext.Flags.DataDirectory
	opts.ValueDir = SystemContext.Flags.DataDirectory
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	node_seq, err := db.GetSequence(NodeIDSequencePath, 10)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	SystemContext.Badger = db
	SystemContext.NodeIDs = node_seq
	fmt.Println("Starting admin application with port:", SystemContext.Flags.HTTPPort)
	fmt.Println("Listening for connections on:", SystemContext.Flags.PostgresPort)
	Database.Start(&SystemContext)
	// api.StartApp(&SystemContext)
}
