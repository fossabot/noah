package main

import (
	"fmt"
	"github.com/Ready-Stock/Noah/Configuration"
	"github.com/Ready-Stock/Noah/Database"
	"github.com/Ready-Stock/Noah/Database/cluster"
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
	s, err := db.GetSequence([]byte("test1"), 1)
	id, err := s.Next()
	fmt.Printf("NEW ID: %d", id)
	Conf.ParseConfiguration()
	fmt.Println("Starting admin application with port:", Conf.Configuration.AdminPort)
	cluster.SetupNodes()
	Database.Start()
}
