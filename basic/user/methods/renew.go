package methods

import (
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
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
	return routing.GenSuccessResponse(u.Balance)
}
