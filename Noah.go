package main

import (
	"fmt"
	"github.com/Ready-Stock/Noah/Configuration"
	"github.com/Ready-Stock/Noah/Database"
	"github.com/Ready-Stock/Noah/Database/cluster"
)

func main() {
	Conf.ParseConfiguration()
	fmt.Println("Starting admin application with port:", Conf.Configuration.AdminPort)
	cluster.SetupNodes()
	Database.Start()
}
