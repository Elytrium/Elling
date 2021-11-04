package main

import (
	"Elling/config"
	"Elling/elling"
	"Elling/routing"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	log.Log().Msg("Starting Elling - Module-based billing platform made with Go")
	log.Log().Msg("(c) 2021 Elytrium")

	config.LoadConfig()
	elling.ReloadModules()
	StartTicker(time.Duration(config.AppConfig.APITick), routing.DoTick)
	StartTicker(time.Duration(config.AppConfig.ModuleSmallTick), elling.DoSmallTick)
	StartTicker(time.Duration(config.AppConfig.ModuleBigTick), elling.DoBigTick)
}

func StartTicker(interval time.Duration, task func()) {
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		for {
			select {
			case <- ticker.C:
				task()
			}
		}
	}()
}
