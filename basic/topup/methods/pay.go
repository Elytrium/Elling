package methods

import (
	"github.com/Elytrium/elling/basic/topup"
	"github.com/Elytrium/elling/elling"
	"github.com/Elytrium/elling/routing"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
)

type Pay struct{}

func (Pay) GetLimit() int {
	return 60
}

func (Pay) GetType() routing.MethodType {
	return routing.Http
}

func (Pay) IsPublic() bool {
	return false
}

func (Pay) Process(u elling.User, p url.Values) routing.HTTPResponse {
	amount, err := strconv.Atoi(p.Get("amount"))

	if err != nil || amount < 1 {
		return routing.GenBadRequestResponse("topup.invalid-amount")
	}

	method := p.Get("method")

	if instruction, ok := topup.Instructions[method]; ok {
		method := instruction.(topup.Method)
		pendingPurchase, err := method.RequestTopUp(u, amount)

		if err != nil {
			log.Err(err)

			return routing.GenInternalServerError("topup.pay-failed")
		}

		return routing.GenSuccessResponse(pendingPurchase.GetPayString())
	}

	return routing.GenBadRequestResponse("topup.no-instruction")
}
