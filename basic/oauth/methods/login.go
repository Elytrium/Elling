package methods

import (
	"github.com/Elytrium/elling/basic/oauth/types"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog/log"
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

func (Login) Process(_ elling.User, v url.Values) routing.HTTPResponse {
	if instruction, ok := types.Instructions[v.Get("name")]; ok {
		service := instruction.(*types.Service)
		linkedAccount, err := service.ToLinkedAccount(v.Get("key"))

		if err != nil {
			log.Trace().Err(err).Msg("OAuth invalid")
			return routing.GenBadRequestResponse("user.oauth-invalid")
		}

		var dbLinkedAccount types.LinkedAccount
		res := elling.DB.Where("id = ?", linkedAccount.ID).Preload("User").First(&dbLinkedAccount)
		if res.Error != nil {
			linkedAccount.User = elling.NewUser()
			log.Trace().Interface("account", linkedAccount).Msg("Creating new user")
			elling.DB.Save(&linkedAccount)
			dbLinkedAccount = linkedAccount
		}

		dbLinkedAccount.DisplayParam = linkedAccount.DisplayParam

		return routing.GenSuccessResponse(dbLinkedAccount)
	}

	return routing.GenBadRequestResponse("oauth.no-instruction")
}
