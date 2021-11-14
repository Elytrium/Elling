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
	log.Log().Msg("(c) 2021 Elytrium")

	elling.LoadDatabase()
	module.ReloadModules()
	routing.InitRouter()

	StartTicker(time.Duration(config.AppConfig.APITick), routing.DoTick)
	StartTicker(time.Duration(config.AppConfig.ModuleSmallTick), module.DoSmallTick)
	StartTicker(time.Duration(config.AppConfig.ModuleBigTick), module.DoBigTick)
}

func StartTicker(interval time.Duration, task func()) {
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				task()
			}
		}
	}()
}
