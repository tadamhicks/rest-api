package config

import (
	"log"
//	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env"
)

type Config struct {
	Server string `env:"RESTING_DB_IP"`
	Database string `env:"RESTING_DB"`
	Username string `env:"RESTING_DB_USER"`
	Password string `env:"RESTING_DB_PASSWORD"`
}


func (c *Config) Read() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse ENV")
    }
}

/*
type Config struct {
	Server string
	Database string
	Username string
	Password string
}
*/

/*
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
*/
