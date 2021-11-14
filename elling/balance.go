package elling

import "time"

type Balance struct {
	ID     int64 `json:"id,omitempty"`
	Amount uint  `json:"amount"`
}

func NewBalance() Balance {
	return Balance{
		ID:     time.Now().UnixNano(),
		Amount: 0,
	}
}
