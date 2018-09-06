package main

import (
	"fmt"
	"github.com/Ready-Stock/Noah/conf"
	"github.com/Ready-Stock/Noah/db"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/badger"
)

func main() {
	opts := badger.DefaultOptions
	opts.Dir = "badge"
	opts.ValueDir = "badge"
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sctx := system.SContext{
		Badger: db,
	}
	conf.ParseConfiguration()
	fmt.Println("Starting admin application with port:", conf.Configuration.AdminPort)
	Database.Start(&sctx)
}
