package elling

import (
	"Elling/routing"
	"fmt"
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

func Load(m Module)  {
	List = append(List, m)
	name := m.GetName()

	m.OnInit()

	tables := m.OnDBMigration()
	for _, table := range tables {
		err := DB.AutoMigrate(table)

		if err != nil {
			log.Err(err)
		}
	}

	methods := m.OnRegisterMethods()
	routing.Router[m.GetName()] = methods

	log.Info().Msg("Loaded module " + name)
}

func ReloadModules()  {
	plugins, err := filepath.Glob("plugins/*.so")

	if err != nil {
		log.Err(err)
	}

	for _, filename := range plugins {
		fmt.Println(filename)
		p, err := plugin.Open(filename)
		if err != nil {
			log.Err(err)
		}

		moduleSym, err := p.Lookup("Module")

		if err != nil {
			log.Err(err)
		}

		module, ok := moduleSym.(Module)

		if !ok {
			log.Err(err)
		}

		Load(module)
	}
}