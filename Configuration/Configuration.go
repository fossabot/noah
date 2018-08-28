package Conf

import (
	"flag"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

var (
	Configuration Config
)

type Config struct {
	AdminPort int            `json:"admin_port"`
	Database  DatabaseConfig `json:"database"`
	Cluster   ClusterConfig  `json:"cluster"`
	Nodes     []NodeConfig   `json:"nodes"`
}

type DatabaseConfig struct {
	AdvertiseAddress string `json:"advertise_address"`
	ReadBuffer       int    `json:"read_buffer"`
}

type ClusterConfig struct {
	DenyConnectionIfNoNodes    bool `json:"deny_connection_if_no_nodes"`
	StartingConnectionsPerNode int  `json:"starting_connections_per_node"`
}

type NodeConfig struct {
	NodeID   int    `json:"node_id"`
	Address  string `json:"address"`
	Port 	 uint16	`json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func ParseConfiguration() {
	configPath := flag.String("config", "config.json", "Path to the Noah config file.")

	if *configPath == "" {
		panic("Error, config path cannot be blank!")
	}

	dat, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic("Could not find config file from (" + string(*configPath) + ")")
	}
	var config Config
	err = json.Unmarshal(dat, &config)
	if err != nil {
		panic("Could not read config file from (" + string(*configPath) + ")")
	}
	Configuration = config
}

func SaveConfiguration() {
	configPath := flag.String("config", "config.json", "Path to the Noah config file.")

	if *configPath == "" {
		panic("Error, config path cannot be blank!")
	}

	if config, err := json.Marshal(Configuration); err != nil {
		fmt.Errorf("Error, could not serialize config!")
	} else if err := ioutil.WriteFile(*configPath, config, os.ModeExclusive); err != nil {
		fmt.Errorf("Error, could not save config!")
	} else {
		fmt.Println("Configuration saved to:", configPath)
	}
}
