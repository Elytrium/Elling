package routing

import (
	"encoding/json"
	"github.com/Elytrium/elling/elling"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

type HTTPResponse struct {
	Response
	Code int
}

type Response struct {
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
	Answer  interface{} `json:"answer,omitempty"`
}

type MethodType int

const (
	Http MethodType = iota
	Socket
)

type Method interface {
	IsPublic() bool
	GetLimit() int
	GetType() MethodType
}

type HTTPMethod interface {
	Method
	Process(user *elling.User, params *url.Values) *HTTPResponse
}

type SocketMethod interface {
	Method
	CanRegister(params *url.Values) *HTTPResponse
	Register(conn *websocket.Conn, user *elling.User)
}

func (r *HTTPResponse) Write(writer http.ResponseWriter) {
	encoder := json.NewEncoder(writer)

	writer.WriteHeader(r.Code)
	_ = encoder.Encode(r.Response)
}

func GenSuccessResponse(answer interface{}) *HTTPResponse {
	return &HTTPResponse{
		Response: Response{
			Success: true,
			Message: "done",
			Answer:  answer,
		},
		Code: 200,
	}
}

func GenBadRequestResponse(message string) *HTTPResponse {
	return &HTTPResponse{
		Response: Response{
			Success: false,
			Message: message,
			Answer:  nil,
		},
		Code: 400,
	}
}

func GenForbiddenResponse(message string) *HTTPResponse {
	return &HTTPResponse{
		Response: Response{
			Success: false,
			Message: message,
			Answer:  nil,
		},
		Code: 403,
	}
}

func GenUnauthorizedResponse(message string) *HTTPResponse {
	return &HTTPResponse{
		Response: Response{
			Success: false,
			Message: message,
			Answer:  nil,
		},
		Code: 401,
	}
}

func GenTooManyRequestsResponse(message string) *HTTPResponse {
	return &HTTPResponse{
		Response: Response{
			Success: false,
			Message: message,
			Answer:  nil,
		},
		Code: 429,
	}
}

func GenInternalServerError(message string) *HTTPResponse {
	return &HTTPResponse{
		Response: Response{
			Success: false,
			Message: message,
			Answer:  nil,
		},
		Code: 500,
	}
}
