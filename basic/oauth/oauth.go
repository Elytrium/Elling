package main

import (
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/basic/oauth/methods"
	"github.com/Elytrium/elling/basic/oauth/types"
	"github.com/Elytrium/elling/module"
	"github.com/Elytrium/elling/routing"
	"reflect"
)

type OAuth struct{}

func (*OAuth) OnModuleInit() {
	types.Instructions = common.ReadInstructions("oauth", reflect.TypeOf(types.Service{}))
}

func (*OAuth) OnModuleRemove() {
	types.Instructions = common.Instructions{}
}

func (*OAuth) GetMeta() *module.Meta {
	return &module.Meta{
		Name: "oauth",
		Routes: map[string]routing.Method{
			"login": &methods.Login{},
			"list":  &methods.List{},
		},
		DatabaseFields: []interface{}{
			&types.LinkedAccount{},
		},
	}
}

var Module OAuth
