package main

import (
	"github.com/BurntSushi/toml"
	"fmt"
)

type config struct {
	DB  database `toml:"database"`
	App app      `toml:"app"`
}

type database struct {
	DBname string
	DBusername string
	DBPassword string
	DBport int
}

type app struct {
	Port int
}

func getConfig() config {
	var cfg config
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		fmt.Println(err)
	}
	return cfg
}
