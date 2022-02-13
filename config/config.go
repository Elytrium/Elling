package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var AppConfig Config

type DBType string

const (
	MySQL      DBType = "mysql"
	PostgreSQL DBType = "postgresql"
	SQLite     DBType = "sqlite"
)

type Config struct {
	APIAddress      string `default:":80"`
	APIRequestLimit int    `default:"256"`
	APITick         int    `default:"60"`
	ModuleSmallTick int    `default:"60"`
	ModuleBigTick   int    `default:"86400"`
	DBType          DBType `default:"sqlite"`
	DSN             string `default:"elling.db"`
	RedisConn       string `default:"127.0.0.1:6379"`
	RedisMaxIdle    int    `default:"3"`
	RedisTimeout    int    `default:"240"`
	RedisPassword   string `default:"redispassword"`
	SlowDBThreshold int64  `default:"400"`
	LogLevel        string `default:"trace"`
	MachineID       string `default:"ip"`
	StartTime       int64  `default:"1640995200000"`
	IsMaster        bool   `default:"false"`
}

func LoadConfig() {
	err := envconfig.Process("elling", &AppConfig)
	if err != nil {
		log.Error().Err(err).Send()
	}
}
