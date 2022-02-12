package main

import (
	"github.com/Elytrium/elling/basic/topup/types"
	"github.com/Elytrium/elling/elling"
	"time"
)

type TopUpListener struct{}

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
