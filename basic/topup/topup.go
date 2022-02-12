package main

import (
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/basic/topup/methods"
	"github.com/Elytrium/elling/basic/topup/types"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"reflect"
)

type TopUp struct{}

func (*TopUp) OnModuleInit() {
	types.Instructions = common.ReadInstructions("topup", reflect.TypeOf(types.Method{}))
	elling.RegisterListener(TopUpListener{})
}

func (*TopUp) OnModuleRemove() {
	types.Instructions = common.Instructions{}
}

var Module TopUp

var ModuleMeta = elling.ModuleMeta{
	Name: "topup",
	Routes: map[string]routing.Method{
		"list": &methods.List{},
		"pay":  &methods.Pay{},
	},
	DatabaseFields: []interface{}{
		types.PendingPurchase{},
	},
}
