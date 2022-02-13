package methods

import (
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Renew struct{}

func (*Renew) GetLimit() int {
	return 10
}

func (*Renew) GetType() routing.MethodType {
	return routing.Http
}

func (*Renew) IsPublic() bool {
	return false
}

func (*Renew) Process(u *elling.User, _ *url.Values) *routing.HTTPResponse {
	err := u.ChangeToken()

	if err != nil {
		log.Error().Err(err).Msg("Changing token")
		return routing.GenInternalServerError("change-token-write")
	}

	return routing.GenSuccessResponse(u)
}
