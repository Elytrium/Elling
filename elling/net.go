package elling

import (
	"encoding/json"
	"github.com/mdaverde/jsonpath"
	"io/ioutil"
	"net/http"
	"strings"
)

type NetRequest struct {
	URL               string
	Method            string
	Headers           map[string]string
	Data              string
	ResponseType      ResponseType
	ResponseValuePath []string
}

func (request NetRequest) DoRequest(replaceValues map[string]string) ([]string, error) {
	client := http.Client{}

	safeRequestURL := request.URL
	safeRequestData := request.Data
	safeHeaders := request.Headers

	for from, to := range replaceValues {
		safeRequestURL = strings.Replace(safeRequestURL, from, to, -1)
	}

	if safeRequestData != "" {
		for from, to := range replaceValues {
			safeRequestData = strings.Replace(safeRequestData, from, to, -1)
		}
	}

	if len(safeHeaders) != 0 {
		for i := range replaceValues {
			for from, to := range replaceValues {
				safeHeaders[i] = strings.Replace(safeHeaders[i], from, to, -1)
			}
		}
	}

	httpRequest, err := http.NewRequest(request.Method, safeRequestURL, strings.NewReader(safeRequestData))

	if err != nil {
		return nil, err
	}

	for headerName, headerValue := range safeHeaders {
		httpRequest.Header.Add(headerName, headerValue)
	}

	httpResponse, err := client.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	switch request.ResponseType {
	case ResponseJson:
		var data interface{}
		var serializedResult []string

		json.NewDecoder(httpResponse.Body).Decode(&data)
		for i := range request.ResponseValuePath {
			responseValuePath := request.ResponseValuePath[i]
			responseValue, err := jsonpath.Get(data, responseValuePath)

			if err != nil {
				return nil, err
			}

			serializedResult = append(serializedResult, responseValue.(string))
		}
		return serializedResult, nil
	case ResponsePlain:
		bytes, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return nil, err
		}

		return []string{string(bytes)}, nil
	}

	return nil, nil
}

type ResponseType string

const (
	ResponseJson  ResponseType = "JSON"
	ResponsePlain              = "PLAIN"
	ResponseNone = "NONE"
)
