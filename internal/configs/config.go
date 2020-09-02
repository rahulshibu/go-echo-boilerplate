package configs

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

//Config - config for application
type Config struct {
	Name   string
	Secret string
	Server struct {
		Host string
		Port string
	}
	Database struct {
		User     string
		Password string
		Host     string
		Port     string
		Name     string
	}
	Environment string
}

// AppConfig is the configs for the whole application
var AppConfig *Config

//Init - initialize config
func Init() error {
	if _, err := toml.DecodeFile("config.toml", &AppConfig); err != nil {
		log.Fatalf(" %s", err)
		return err
	}
	return nil
}
