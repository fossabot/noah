package main

import (
	"flag"
	"fmt"
	"github.com/Ready-Stock/Noah/db"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/badger"
	"github.com/fatih/color"
)

var (
	argHTTPPort      = flag.Int("http-port", 8080, "Listen port for Noah's web interface.")
	argPostgresPort  = flag.Int("psql-port", 5433, "Listen port for Noah's PostgreSQL client connectivity.")
	argDataDirectory = flag.String("data-dir", "badge", "Directory for Noah's embedded database.")
)

func main() {
	validateFlags()
	opts := badger.DefaultOptions
	opts.Dir = *argDataDirectory
	opts.ValueDir = *argDataDirectory
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sctx := system.SContext{
		Badger: db,
	}
	fmt.Println("Starting admin application with port:", *argHTTPPort)
	Database.Start(&sctx)
}

func validateFlags() {
	if *argHTTPPort < 1 {
		color.Red("Error, HTTP port cannot be less than 1. Invalid port provided: %d", argHTTPPort)
	}
}
