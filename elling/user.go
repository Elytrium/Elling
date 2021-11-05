package elling

import "github.com/Elytrium/elling/utils"

type User struct {
	DBModel
	Email         string        `json:"email,omitempty"`
	Balance       Balance       `json:"balance"`
	Products      []Product     `json:"products,omitempty"`
	LinkedAccount LinkedAccount `json:"linked_accounts,omitempty"`
	Hash          string        `json:"hash,omitempty"`
	Token         string        `json:"token,omitempty"`
}

type LinkedAccount struct {
	ID           string `json:"id,omitempty"`
	DisplayParam string `json:"display_param,omitempty"`
}

type OAuthService struct {
	DisplayName     string     `yaml:"display-name" json:"display_name,omitempty"`
	Name            string     `yaml:"name" json:"name,omitempty"`
	OAuthGenRequest string     `yaml:"oauth-gen-request" json:"oauth_gen_request,omitempty"`
	VerifyRequest   NetRequest `yaml:"verify-request" json:"verify_request"`
}

func (t OAuthService) ToLinkedAccount(token string) (LinkedAccount, error) {
	resp, err := t.VerifyRequest.DoRequest(map[string]string{
		"{token}": token,
	})

	if err != nil {
		return LinkedAccount{}, err
	}

	displayParam := resp[0]
	id := resp[1]

	return LinkedAccount{
		ID:           id,
		DisplayParam: displayParam,
	}, nil
}

func NewUser(email string, account LinkedAccount) User {
	user := User{
		Email:         email,
		LinkedAccount: account,
		Balance:       NewBalance(),
		Token:         utils.GenToken(64),
	}

	return user
}
