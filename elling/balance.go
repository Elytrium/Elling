package elling

type Balance struct {
	ID     uint64 `json:"id,omitempty"`
	Amount int64  `json:"amount"`
}

type BalanceChangeEvent struct {
	User  *User
	Delta int64
}

func NewBalance() *Balance {
	return &Balance{
		ID:     NextID(),
		Amount: 0,
	}
}
