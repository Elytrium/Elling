package elling

import "errors"

type Balance struct {
	ID     uint64 `json:"id,omitempty"`
	Amount int64  `json:"amount"`
}

type BalanceChangeEvent struct {
	Balance *Balance
	Delta   int64
}

var NotEnoughFundsError = errors.New("not enough funds")

func NewBalance() (*Balance, error) {
	balance := &Balance{
		ID:     NextID(),
		Amount: 0,
	}

	return balance, balance.Update()
}

func (b *Balance) Deposit(delta int64) error {
	b.Amount += delta
	DispatchEvent(BalanceChangeEvent{Balance: b, Delta: delta})
	return b.Update()
}

func (b *Balance) Withdraw(delta int64) error {
	b.Amount -= delta
	DispatchEvent(BalanceChangeEvent{Balance: b, Delta: -delta})
	return b.Update()
}

func (b *Balance) Update() error {
	DB.Save(b)
	return DB.Error
}
