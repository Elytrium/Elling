package elling

import (
	"encoding/json"
	"errors"
	"github.com/mdaverde/jsonpath"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type NetRequest struct {
	URL               string            `yaml:"url"`
	Method            string            `yaml:"method"`
	Headers           map[string]string `yaml:"headers"`
	Data              string            `yaml:"data"`
	ResponseType      ResponseType      `yaml:"response-type"`
	ResponseValuePath []string          `yaml:"response-value-path"`
}

var PathNotFoundError = errors.New("url decoding path not found")

func (request *NetRequest) DoRequest(replaceValues map[string]string) ([]string, error) {
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
		for i := range safeHeaders {
			for from, to := range replaceValues {
				safeHeaders[i] = strings.Replace(safeHeaders[i], from, to, -1)
			}
		}
	}

	log.Trace().Interface("headers", safeHeaders).Str("data", safeRequestData).Str("url", safeRequestURL).Str("method", request.Method).Msg("Sending request")

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

		err := json.NewDecoder(httpResponse.Body).Decode(&data)

		if err != nil {
			return nil, err
		}

		for _, responseValuePath := range request.ResponseValuePath {
			responseValue, err := jsonpath.Get(data, responseValuePath)

			if err != nil {
				return nil, err
			}

			var value string
			value, ok := responseValue.(string)
			if !ok {
				value = strconv.Itoa(int(responseValue.(float64)))
			}

			serializedResult = append(serializedResult, value)
		}
		return serializedResult, nil
	case ResponseURLEncoded:
		var serializedResult []string
		bytes, err := ioutil.ReadAll(httpResponse.Body)

		params, err := url.ParseQuery(string(bytes))
		if err != nil {
			return nil, err
		}
		for _, responseValuePath := range request.ResponseValuePath {
			responseValue, ok := params[responseValuePath]

			if !ok {
				return nil, PathNotFoundError
			}

			serializedResult = append(serializedResult, responseValue[0])
		}
		return serializedResult, nil
	case ResponsePlain:
		bytes, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return nil, err
		}

		return []string{string(bytes)}, nil
	case ResponseNone:
		return nil, nil
	default:
		return nil, errors.New("unknown response type " + string(request.ResponseType))
	}
}

type ResponseType string

const (
	ResponseJson       ResponseType = "JSON"
	ResponseURLEncoded              = "URLEncoded"
	ResponsePlain                   = "PLAIN"
	ResponseNone                    = "NONE"
)
