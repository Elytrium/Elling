package main

import (
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/basic/oauth/methods"
	"github.com/Elytrium/elling/basic/oauth/types"
	"github.com/Elytrium/elling/routing"
	"reflect"
)

type OAuth struct{}

func (o OAuth) OnInit() {
	types.Instructions = common.ReadInstructions("oauth", reflect.TypeOf(types.Service{}))
}

func (o OAuth) GetName() string {
	return "oauth"
}

func (o OAuth) OnRegisterMethods() map[string]routing.Method {
	return map[string]routing.Method{
		"login": methods.Login{},
		"list":  methods.List{},
	}
}

func (o OAuth) OnDBMigration() []interface{} {
	return []interface{}{
		&types.LinkedAccount{},
	}
}

func (o OAuth) OnSmallTick() {

}

func (o OAuth) OnBigTick() {

}

var Module OAuth
