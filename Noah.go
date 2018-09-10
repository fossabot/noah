package main

import (
	"flag"
	"fmt"
	"github.com/Ready-Stock/Noah/db"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/badger"
)

func main() {
	SystemContext := system.SContext{
		Flags: system.SFlags{
			HTTPPort:      *flag.Int("http-port", 8080, "Listen port for Noah's web interface."),
			PostgresPort:  *flag.Int("psql-port", 5433, "Listen port for Noah's PostgreSQL client connectivity."),
			DataDirectory: * flag.String("data-dir", "badge", "Directory for Noah's embedded database."),
		},
	}
	opts := badger.DefaultOptions
	opts.Dir = SystemContext.Flags.DataDirectory
	opts.ValueDir = SystemContext.Flags.DataDirectory
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	SystemContext.Badger = db
	fmt.Println("Starting admin application with port:", SystemContext.Flags.HTTPPort)
	fmt.Println("Listening for connections on:", SystemContext.Flags.PostgresPort)
	Database.Start(&SystemContext)
}
