package module

import (
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"plugin"
)

type Module interface {
	OnInit()
	GetName() string
	OnRegisterMethods() map[string]routing.Method
	OnDBMigration() []interface{}
	OnSmallTick()
	OnBigTick()
}

var List []Module

func DoSmallTick() {
	for _, module := range List {
		module.OnSmallTick()
	}
}

func DoBigTick() {
	for _, module := range List {
		module.OnBigTick()
	}
}

func Load(m Module) {
	List = append(List, m)
	name := m.GetName()

	m.OnInit()

	tables := m.OnDBMigration()
	for _, table := range tables {
		err := elling.DB.AutoMigrate(table)

		if err != nil {
			log.Error().Err(err).Send()
		}
	}

	methods := m.OnRegisterMethods()
	routing.Router[m.GetName()] = methods

	log.Info().Str("name", name).Msg("Loaded module " + name)
}

func ReloadModules() {
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
			log.Error().Err(err).Send()
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
