package methods

import (
	"github.com/Elytrium/elling/basic/oauth"
	"github.com/Elytrium/elling/routing"
)

type List struct{}

func (List) GetLimit() int {
	return 60
}

func (List) GetType() routing.MethodType {
	return routing.Http
}

func (List) IsPublic() bool {
	return true
}

func (List) Process() routing.HTTPResponse {
	return routing.GenSuccessResponse(oauth.Instructions)
}
