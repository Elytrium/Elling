package main

import (
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/basic/topup/methods"
	"github.com/Elytrium/elling/basic/topup/types"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/module"
	"github.com/Elytrium/elling/routing"
	"reflect"
	"time"
)

type TopUp struct{}

type TopUpListener struct{}

func (*TopUp) OnModuleInit() {
	types.Instructions = common.ReadInstructions("topup", reflect.TypeOf(types.Method{}))
	elling.RegisterListener(TopUpListener{})
}

func (*TopUp) OnModuleRemove() {
	types.Instructions = common.Instructions{}
}

func (*TopUpListener) OnSmallTick(_ elling.SmallTickEvent) {
	var invalidPurchases []types.PendingPurchase
	elling.DB.Where("InvalidationDate > ?", time.Now()).Find(&invalidPurchases)
	for _, purchase := range invalidPurchases {
		purchase.Reject()
	}

	var validPurchases []types.PendingPurchase
	elling.DB.Find(&validPurchases)
	for _, purchase := range validPurchases {
		purchase.Validate()
	}
}

func (*TopUp) GetMeta() *module.Meta {
	return &module.Meta{
		Name: "topup",
		Routes: map[string]routing.Method{
			"list": &methods.List{},
			"pay":  &methods.Pay{},
		},
		DatabaseFields: []interface{}{
			types.PendingPurchase{},
		},
	}
}

var Module TopUp
