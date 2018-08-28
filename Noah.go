package main

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db"
	"github.com/Ready-Stock/badger"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/Noah/conf"
)

func main() {
	opts := badger.DefaultOptions
	opts.Dir = ""
	opts.ValueDir = ""
	badge, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	sctx := system.SContext{
		Badger: badge,
	}
	defer badge.Close()
	conf.ParseConfiguration()
	if err := sctx.UpsertNodes(conf.Configuration.Nodes); err != nil {
		panic(err)
	}
	if err := sctx.StartConnectionPool(); err != nil {
		panic(err)
	}
	fmt.Println("Starting admin application with port:", conf.Configuration.AdminPort)
	Database.Start(&sctx)
}
