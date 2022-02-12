package methods

import (
	"github.com/Elytrium/elling/basic/topup/types"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"net/url"
)

type List struct{}

func (*List) GetLimit() int {
	return 60
}

func (*List) GetType() routing.MethodType {
	return routing.Http
}

func (*List) IsPublic() bool {
	return true
}

func (*List) Process(_ *elling.User, _ *url.Values) *routing.HTTPResponse {
	return routing.GenSuccessResponse(types.Instructions)
}
