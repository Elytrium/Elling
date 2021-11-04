package elling

type Balance struct {
	ID     int64
	Amount uint
}

func NewBalance() Balance {
	return Balance{}
}
