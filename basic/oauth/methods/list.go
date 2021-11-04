package methods

import (
	"Elling/basic/oauth"
	"Elling/routing"
)

type List struct {}

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