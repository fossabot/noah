package main

import (
	"fmt"
	"github.com/Ready-Stock/Noah/Configuration"
	"github.com/Ready-Stock/Noah/Database"
)

func main() {
	Conf.ParseConfiguration()
	fmt.Println("Starting admin application with port:", Conf.Configuration.AdminPort)
	Database.SetupNodes()
	Database.Start()
}
