package module

import (
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"plugin"
	"reflect"
)

var Modules []Module

type Module interface {
	GetMeta() *Meta
	OnModuleInit()
	OnModuleRemove()
}

type Meta struct {
	Name           string
	Routes         map[string]routing.Method
	DatabaseFields []interface{}
}

func Load(module Module) {
	Modules = append(Modules, module)

	meta := module.GetMeta()
	module.OnModuleInit()

	for _, table := range meta.DatabaseFields {
		err := elling.DB.AutoMigrate(table)

		if err != nil {
			log.Error().Err(err).Send()
		}
	}

	routing.AddRoute(meta.Name, meta.Routes)
	log.Info().Str("name", meta.Name).Msg("Loaded module " + meta.Name)
}

func LoadModules() {
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

		module, ok := moduleSym.(Module)

		if !ok {
			log.Error().Str("filename", filename).Msg("Error in loading plugin: wrong Module format")
			continue
		}

		Load(module)
	}
}

func UnloadPlugins() {
	for _, module := range Modules {
		module.OnModuleRemove()
	}

	elling.ModuleDispatchers = make(map[reflect.Type][]*elling.Dispatcher)

	Modules = []Module{}
}
