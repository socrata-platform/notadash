package lib

import (
	"code.google.com/p/gcfg"
	"log"
)

type Config struct {
	Environment map[string]*struct {
		Marathon_Host string
		Mesos_Host    string
	}
}

func LoadCfg(file string) Config {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, file)
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}
	return cfg
}
