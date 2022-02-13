package main

import (
	"github.com/Elytrium/elling/basic/user/methods"
	"github.com/Elytrium/elling/module"
	"github.com/Elytrium/elling/routing"
)

type User struct{}

func (*User) OnModuleInit() {}

func (*User) OnModuleRemove() {}

func (*User) GetMeta() *module.Meta {
	return &module.Meta{
		Name: "user",
		Routes: map[string]routing.Method{
			"info":        &methods.Info{},
			"renew_token": &methods.Renew{},
		},
		DatabaseFields: []interface{}{},
	}
}

var Module User
