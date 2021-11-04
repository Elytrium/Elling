package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var AppConfig Config

type Config struct {
	APIAddress      string `default:":80"`
	APIRequestLimit int    `default:"256"`
	APITick         int    `default:"60"`
	ModuleSmallTick int    `default:"60"`
	ModuleBigTick   int    `default:"86400"`
}

func LoadConfig() {
	err := envconfig.Process("elling", &AppConfig)
	if err != nil {
		log.Err(err)
	}
}
