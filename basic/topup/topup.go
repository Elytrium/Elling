package topup

import (
	"Elling/basic/common"
	"Elling/basic/topup/methods"
	"Elling/elling"
	"Elling/routing"
	"reflect"
	"time"
)

type TopUp struct {}

var Instructions common.Instructions

func (o TopUp) OnInit() {
	Instructions = common.ReadInstructions("topup", reflect.TypeOf(Method{}))
}

func (o TopUp) GetName() string {
	return "oauth"
}

func (o TopUp) OnRegisterMethods() map[string]routing.Method {
	return map[string]routing.Method{
		"list": methods.List{},
		"pay": methods.Pay{},
	}
}

func (o TopUp) OnDBMigration() []interface{} {
	return []interface{} {
		PendingPurchase{},
	}
}

func (o TopUp) OnSmallTick() {
	var invalidPurchases []PendingPurchase
	elling.DB.Where("InvalidationDate > ?", time.Now()).Find(&invalidPurchases)
	for _, purchase := range invalidPurchases {
		purchase.Reject()
	}

	var validPurchases []PendingPurchase
	elling.DB.Find(&validPurchases)
	for _, purchase := range validPurchases {
		purchase.Validate()
	}
}

func (o TopUp) OnBigTick() {

}

var Module TopUp