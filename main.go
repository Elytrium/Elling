package main

import (
	"github.com/Elytrium/elling/config"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/module"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	config.LoadConfig()

	logLevel, err := zerolog.ParseLevel(config.AppConfig.LogLevel)
	log.Err(err).Msg("Initializing logger")
	zerolog.SetGlobalLevel(logLevel)

	log.Log().Msg("Starting Elling - Module-based billing platform made with Go")
	log.Log().Msg("(c) 2021-2022 Elytrium")

	elling.InitID()
	elling.LoadDatabase()
	module.ReloadModules()

	StartTicker(time.Duration(config.AppConfig.APITick)*time.Second, routing.DoTick)
	StartTicker(time.Duration(config.AppConfig.ModuleSmallTick)*time.Second, elling.DoSmallTick)
	StartTicker(time.Duration(config.AppConfig.ModuleBigTick)*time.Second, elling.DoBigTick)

	routing.InitRouter()
}

func StartTicker(interval time.Duration, task func()) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				task()
			}
		}
	}()
}
