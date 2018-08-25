package main

import (
	"fmt"
	"github.com/Ready-Stock/Noah/Configuration"
	"github.com/Ready-Stock/Noah/Database"
	"github.com/Ready-Stock/Noah/Database/cluster"
	"github.com/Ready-Stock/badger"
	"github.com/Ready-Stock/Noah/Database/system"
)

func main() {
	opts := badger.DefaultOptions
	opts.Dir = "badge"
	opts.ValueDir = "badge"
	badge, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	sctx := system.SContext{
		Badger: badge,
	}
	defer badge.Close()
	Conf.ParseConfiguration()
	fmt.Println("Starting admin application with port:", Conf.Configuration.AdminPort)
	cluster.SetupNodes()
	Database.Start(&sctx)
}
