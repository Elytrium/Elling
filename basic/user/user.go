package main

import (
	"github.com/Elytrium/elling/basic/user/methods"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
)

type User struct{}

func (*User) OnModuleInit() {}

func (*User) OnModuleRemove() {}

var Module User

var ModuleMeta = elling.ModuleMeta{
	Name: "user",
	Routes: map[string]routing.Method{
		"info":        &methods.Info{},
		"renew_token": &methods.Renew{},
	},
	DatabaseFields: []interface{}{},
}
