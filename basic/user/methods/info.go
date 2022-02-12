package methods

import (
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Info struct{}

func (*Info) GetLimit() int {
	return 60
}

func (*Info) GetType() routing.MethodType {
	return routing.Http
}

func (*Info) IsPublic() bool {
	return false
}

func (*Info) Process(u *elling.User, _ *url.Values) *routing.HTTPResponse {
	err := u.FindProducts()

	if err != nil {
		log.Error().Err(err).Msg("Find Association")
		return routing.GenInternalServerError("find-association")
	}

	return routing.GenSuccessResponse(u)
}
