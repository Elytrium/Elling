package elling

import (
	"github.com/Elytrium/elling/utils"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	BalanceID int64     `json:"-"`
	Balance   Balance   `json:"balance"`
	Products  []Product `json:"products,omitempty"`
	Token     string    `json:"token,omitempty"`
}

func NewUser() User {
	user := User{
		ID:       time.Now().UnixNano(),
		Balance:  NewBalance(),
		Products: []Product{},
		Token:    utils.GenToken(64),
	}

	return user
}
