package main

var (
	Configuration Config
)

type Config struct {
	AdminPort int `yaml:"admin_port"`
	DataPort  int `yaml:"data_port"`
}
