package elling

import (
	"github.com/Elytrium/elling/utils"
)

type User struct {
	ID        uint64    `json:"id"`
	BalanceID int64     `json:"-"`
	Balance   Balance   `json:"balance"`
	Products  []Product `json:"products,omitempty"`
	Token     string    `json:"token,omitempty"`
	Active    bool      `json:"active"`
}

func NewUser() *User {
	user := &User{
		ID:       NextID(),
		Balance:  *NewBalance(),
		Products: []Product{},
		Token:    utils.GenToken(64),
		Active:   true,
	}

	return user
}

type UserCreationEvent struct {
	User *User
}

type UserActivationEvent struct {
	User *User
}

type UserDeactivationEvent struct {
	User *User
}

func (u *User) FindProducts() error {
	return DB.Model(&u).Association("Products").Find(&u.Products)
}
