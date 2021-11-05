package methods

import (
	"github.com/Elytrium/elling/basic/oauth"
	"github.com/Elytrium/elling/routing"
	"net/url"
)

type Login struct{}

func (Login) GetLimit() int {
	return 10
}

func (Login) GetType() routing.MethodType {
	return routing.Http
}

func (Login) IsPublic() bool {
	return true
}

func (Login) Process(v url.Values) routing.HTTPResponse {
	if instruction, ok := oauth.Instructions[v.Get("name")]; ok {
		service := instruction.(oauth.Service)
		linkedAccount, err := service.ToLinkedAccount(v.Get("key"))

		if err != nil {
			return routing.GenBadRequestResponse("user.oauth-invalid")
		}

		return routing.GenSuccessResponse(linkedAccount)
	}

	return routing.GenBadRequestResponse("oauth.no-instruction")
}
