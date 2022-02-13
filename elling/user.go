package elling

import (
	"github.com/Elytrium/elling/utils"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID        uint64    `json:"id"`
	BalanceID int64     `json:"-"`
	Balance   Balance   `json:"balance"`
	Products  []Product `json:"products,omitempty" gorm:"many2many:user_products;"`
	Token     string    `json:"token,omitempty"`
	Active    bool      `json:"active"`
}

func NewUser() *User {
	bal, err := NewBalance()

	if err != nil {
		log.Error().Err(err).Msg("Creating balance")
	}

	user := &User{
		ID:       NextID(),
		Balance:  *bal,
		Products: []Product{},
		Token:    utils.GenToken(64),
		Active:   true,
	}

	DispatchEvent(UserCreationEvent{User: user})

	DB.Save(user)

	return user
}

func (u *User) Activate() {
	u.Active = true
	DispatchEvent(UserActivationEvent{User: u})
}

func (u *User) Deactivate() {
	u.Active = false
	DispatchEvent(UserDeactivationEvent{User: u})
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

func (u *User) FetchProducts() error {
	return DB.Model(&u).Association("Products").Find(&u.Products)
}
