package module

import (
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"plugin"
)

var Modules []Module

type Module interface {
	OnModuleInit()
	OnModuleRemove()
}

func Load(module Module, meta elling.ModuleMeta) {
	Modules = append(Modules, module)

	module.OnModuleInit()

	for _, table := range meta.DatabaseFields {
		err := elling.DB.AutoMigrate(table)

		if err != nil {
			log.Error().Err(err).Send()
		}
	}

	routing.Router[meta.Name] = meta.Routes
	log.Info().Str("name", meta.Name).Msg("Loaded module " + meta.Name)
}

func ReloadModules() {
	for _, module := range Modules {
		module.OnModuleRemove()
	}

	Modules = []Module{}

	plugins, err := filepath.Glob("plugins/*.so")

	if err != nil {
		log.Error().Err(err).Send()
	}

	for _, filename := range plugins {
		log.Info().Str("filename", filename).Msg("Loading plugin " + filename)
		p, err := plugin.Open(filename)

		if err != nil {
			log.Error().Err(err).Send()
			continue
		}

		moduleSym, err := p.Lookup("Module")

		if err != nil {
			log.Error().Err(err).Msg("Error in loading plugin: Module not found")
			continue
		}

		metaSym, err := p.Lookup("ModuleMeta")

		if err != nil {
			log.Error().Err(err).Msg("Error in loading plugin: ModuleMeta not found")
			continue
		}

		module, ok := moduleSym.(Module)

		if !ok {
			log.Error().Str("filename", filename).Msg("Error in loading plugin: wrong Module format")
			continue
		}

		meta, ok := metaSym.(elling.ModuleMeta)

		if !ok {
			log.Error().Str("filename", filename).Msg("Error in loading plugin: wrong ModuleMeta format")
			continue
		}

		Load(module, meta)
	}
}
