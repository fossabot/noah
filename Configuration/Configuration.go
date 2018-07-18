package Conf

import (
	"flag"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var (
	Configuration Config
)

type Config struct {
	AdminPort int `yaml:"admin_port"`
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	AdvertiseAddress string `yaml:"advertise_address"`
	ReadBuffer int `yaml:"read_buffer"`
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