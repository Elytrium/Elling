package oauth

import (
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/basic/oauth/methods"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"reflect"
)

type OAuth struct{}

var Instructions common.Instructions

func (o OAuth) OnInit() {
	Instructions = common.ReadInstructions("oauth", reflect.TypeOf(Service{}))
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
		elling.LinkedAccount{},
	}
}

func (o OAuth) OnSmallTick() {

}

func (o OAuth) OnBigTick() {

}

var Module OAuth
