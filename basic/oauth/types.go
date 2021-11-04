package oauth

import "Elling/elling"

type Service struct {
	DisplayName     string            `yaml:"display-name" json:"display_name,omitempty"`
	Name            string            `yaml:"name" json:"name,omitempty"`
	OAuthGenRequest string            `yaml:"oauth-gen-request" json:"oauth_gen_request,omitempty"`
	VerifyRequest   elling.NetRequest `yaml:"verify-request" json:"verify_request"`
}

func (t Service) ToLinkedAccount(token string) (elling.LinkedAccount, error) {
	resp, err := t.VerifyRequest.DoRequest(map[string]string{
		"{token}": token,
	})

	if err != nil {
		return elling.LinkedAccount{}, err
	}

	displayParam := resp[0]
	id := resp[1]

	return elling.LinkedAccount{
		ID:           id,
		DisplayParam: displayParam,
	}, nil
}
