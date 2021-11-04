package routing

import (
	"Elling/config"
	"Elling/elling"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var limitAllMap = make(map[string]int)
var limitMap = make(map[string]map[Method]int)

var Router Routes

type NotFoundHandler struct {}

type MethodNotAllowedHandler struct {}

type Routes map[string]map[string]Method

func InitRouter() {
	log.Debug().Interface("routes", Router).Msg("Initializing router")

	router := mux.NewRouter()

	router.HandleFunc("/{group}/{method}", func(writer http.ResponseWriter, request *http.Request) {
		HandleAPI(writer, request)
	})

	router.NotFoundHandler = NotFoundHandler{}
	router.MethodNotAllowedHandler = MethodNotAllowedHandler{}

	address := config.AppConfig.APIAddress
	_ = http.ListenAndServe(address, router)
	log.Info().Str("address", address).Msg("HTTP Listener started")
}

func HandleAPI(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	query := request.URL.Query()

	writer.Header().Add("content-type", "application/json")

	if curLimit, ok := limitAllMap[request.RemoteAddr]; ok {
		if curLimit > config.AppConfig.APIRequestLimit {
			GenTooManyRequestsResponse("all-limit").Write(writer)
			return
		}

		limitAllMap[request.RemoteAddr] += 1
	} else {
		limitAllMap[request.RemoteAddr] = 1
		limitMap[request.RemoteAddr] = make(map[Method]int)
	}

	group := vars["group"]
	method := vars["method"]

	if curGroup, ok := Router[group]; ok {
		if curMethod, ok := curGroup[method]; ok {
			if curLimit, ok := limitMap[request.RemoteAddr][curMethod]; ok {
				if curLimit > curMethod.GetLimit() {
					GenTooManyRequestsResponse(group + "-" + method + "-limit").Write(writer)
					return
				}

				limitMap[request.RemoteAddr][curMethod] += 1
			} else {
				limitMap[request.RemoteAddr][curMethod] = 1
			}

			var userModel elling.User
			if !curMethod.IsPublic() {
				token := request.Header.Get("authorization")
				if token == "" {
					GenForbiddenResponse("private-method").Write(writer)
					return
				}

				q := elling.DB.First(&userModel, "token = ?", token)

				if q.Error != nil {
					GenUnauthorizedResponse("invalid-token").Write(writer)
					return
				}
			}

			switch curMethod.GetType() {
			case Http:
				method := curMethod.(HTTPMethod)
				result := method.Process(userModel, query)

				result.Write(writer)
				break
			case Socket:
				method := curMethod.(SocketMethod)
				result := method.CanRegister(query)

				if result.Success {
					conn, _ := upgrader.Upgrade(writer, request, nil)
					method.Register(conn, userModel)
				} else {
					result.Write(writer)
				}
				break
			}

			return
		} else {
			GenBadRequestResponse("group-method").Write(writer)
		}
	} else {
		GenBadRequestResponse("group").Write(writer)
	}
}

func DoTick() {
	limitAllMap = make(map[string]int)
	limitMap = make(map[string]map[Method]int)
}

func (h NotFoundHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	GenBadRequestResponse("not-found").Write(writer)
}

func (h MethodNotAllowedHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	GenBadRequestResponse("not-allowed").Write(writer)
}
