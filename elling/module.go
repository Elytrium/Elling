package elling

import "github.com/Elytrium/elling/routing"

type ModuleMeta struct {
	Name           string
	Routes         map[string]routing.Method
	DatabaseFields []interface{}
}
