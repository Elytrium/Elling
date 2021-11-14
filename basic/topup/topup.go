package main

import (
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/basic/topup/methods"
	"github.com/Elytrium/elling/basic/topup/types"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"reflect"
	"time"
)

type TopUp struct{}

func (o TopUp) OnInit() {
	types.Instructions = common.ReadInstructions("topup", reflect.TypeOf(types.Method{}))
}

func (o TopUp) GetName() string {
	return "topup"
}

func (o TopUp) OnRegisterMethods() map[string]routing.Method {
	return map[string]routing.Method{
		"list": methods.List{},
		"pay":  methods.Pay{},
	}
}

func (o TopUp) OnDBMigration() []interface{} {
	return []interface{}{
		types.PendingPurchase{},
	}
}

func (o TopUp) OnSmallTick() {
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

func (o TopUp) OnBigTick() {

}

var Module TopUp
