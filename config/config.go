package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
)

/*
type Config struct {
	Server   string `envconfig:"MONGO_IP"`
	Port     string `envconfig:"MONGO_PORT"`
	Database string `envconfig:"MONGO_DB"`
	Username string `envconfig:"MONGO_DB_USER"`
	Password string `envconfig:"MONGO_DB_PASSWORD"`
}
*/

type Config struct {
	Server           string `env:"MONGO_IP"`
	Port             string `env:"MONGO_PORT"`
	Database         string `env:"RESTING_DB"`
	Username         string `env:"RESTING_DB_USER"`
	Password         string `env:"RESTING_DB_PASSWORD"`
	Honeycombkey     string `env:"HONEYCOMB_API_KEY"`
	Honeycombdataset string `env:"HONEYCOMB_DATASET"`
	Servicename      string `env:"SERVICE_NAME"`
}

/*
func (c *Config) Read() {
	err := envconfig.Process("mongo", &c)
	if err != nil {
		log.Fatal(err.Error())
		//log.Fatalf("Failed to parse ENV")
	}
	fmt.Println("env vars read")
	fmt.Printf("%+v\n", c)
}
*/

func (cfg *Config) Read() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
		//log.Fatalf("Failed to parse ENV")
	}
	fmt.Println(cfg)
}

/*
type Config struct {
	Server   string
	Port     string
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
	fmt.Printf("%+v\n", c)
}
*/
