package types

import (
	"errors"
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/elling"
)

var Instructions common.Instructions

type Service struct {
	DisplayName     string            `yaml:"display-name" json:"display_name,omitempty"`
	Name            string            `yaml:"name" json:"name,omitempty"`
	OAuthGenRequest string            `yaml:"oauth-gen-request" json:"oauth_gen_request,omitempty"`
	NeedVerify      bool              `yaml:"need-verify" json:"-"`
	VerifyRequest   elling.NetRequest `yaml:"verify-request" json:"-"`
	GetDataRequest  elling.NetRequest `yaml:"get-data-request" json:"-"`
}

type LinkedAccount struct {
	ID           string      `json:"id,omitempty"`
	UserID       int64       `json:"-"`
	User         elling.User `json:"user"`
	DisplayParam string      `json:"display_param,omitempty" gorm:"-"`
}

var TooShortVerifyAnswer = errors.New("too short verify answer")

var TooShortDataAnswer = errors.New("too short data answer")

func (t Service) ToLinkedAccount(token string) (*LinkedAccount, error) {
	if t.NeedVerify {
		resp, err := t.VerifyRequest.DoRequest(map[string]string{
			"{token}": token,
		})

		if err != nil {
			return nil, err
		}

		if len(resp) < 1 {
			return nil, TooShortVerifyAnswer
		}

		token = resp[0]
	}

	resp, err := t.GetDataRequest.DoRequest(map[string]string{
		"{token}": token,
	})

	if err != nil {
		return nil, err
	}

	if len(resp) < 2 {
		return nil, TooShortDataAnswer
	}

	displayParam := resp[0]
	id := resp[1]

	return &LinkedAccount{
		ID:           id,
		DisplayParam: displayParam,
	}, nil
}
