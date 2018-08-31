package main

import (
	"fmt"
	"github.com/Ready-Stock/Noah/conf"
	"github.com/Ready-Stock/Noah/db"
	"github.com/Ready-Stock/Noah/db/system"
)

func main() {
	sctx := system.SContext{
		Badger: nil,
	}
	conf.ParseConfiguration()
	fmt.Println("Starting admin application with port:", conf.Configuration.AdminPort)
	Database.Start(&sctx)
}
