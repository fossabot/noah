package Conf

import (
	"flag"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"os"
)

var (
	Configuration Config
)

type Config struct {
	AdminPort int            `yaml:"admin_port"`
	Database  DatabaseConfig `yaml:"database"`
	Nodes     []NodeConfig   `yaml:"nodes"`
}

type DatabaseConfig struct {
	AdvertiseAddress string `yaml:"advertise_address"`
	ReadBuffer       int    `yaml:"read_buffer"`
}

type NodeConfig struct {
	NodeID   int    `yaml:"node_id"`
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Shards   int64  `yaml:"shards"`
}

func ParseConfiguration() {
	configPath := flag.String("config", "config.yaml", "Path to the Noah config file.")

	if *configPath == "" {
		panic("Error, config path cannot be blank!")
	}

	dat, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic("Could not find config file from (" + string(*configPath) + ")")
	}
	var config Config
	err = yaml.Unmarshal(dat, &config)
	if err != nil {
		panic("Could not read config file from (" + string(*configPath) + ")")
	}
	Configuration = config
}

func SaveConfiguration() {
	configPath := flag.String("config", "config.yaml", "Path to the Noah config file.")

	if *configPath == "" {
		panic("Error, config path cannot be blank!")
	}

	if config, err := yaml.Marshal(Configuration); err != nil {
		fmt.Errorf("Error, could not serialize config!")
	} else if err := ioutil.WriteFile(*configPath, config, os.ModeExclusive); err != nil {
		fmt.Errorf("Error, could not save config!")
	} else {
		fmt.Println("Configuration saved to:", configPath)
	}
}
