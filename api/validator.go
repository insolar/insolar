// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/insolar/rpc/v2/json2"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type RequestValidator struct {
	router      *openapi3filter.Router
	nextHandler http.Handler
}

func NewRequestValidator(path string, next http.Handler) (*RequestValidator, error) {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load swagger file")
	}
	swagger.Servers = nil

	router := openapi3filter.NewRouter()
	err = router.AddSwagger(swagger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add swagger to router")
	}

	return &RequestValidator{
		router:      router,
		nextHandler: next,
	}, nil
}

type jsonRPCErrorResponse struct {
	Version string `json:"jsonrpc"`

	Error *json2.Error `json:"error,omitempty"`

	ID *json.RawMessage `json:"id"`
}

func (rv *RequestValidator) ServeHTTP(w http.ResponseWriter, httpReq *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			inslogger.FromContext(httpReq.Context()).Error(errors.Errorf("panic in validator: %v", r))
		}
	}()
	err := rv.Validate(httpReq.Context(), httpReq)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(jsonRPCErrorResponse{
			Version: "2.0",
			Error: &json2.Error{
				Code:    InvalidRequestError,
				Message: InvalidRequestErrorMessage,
				Data: requester.Data{
					Trace: []string{err.Error()},
				},
			},
			ID: nil,
		})
		if err != nil {
			inslogger.FromContext(httpReq.Context()).Panic(errors.Wrap(err, "failed to encode error message"))
		}

		return
	}

	rv.nextHandler.ServeHTTP(w, httpReq)
}

var (
	findMethodNameRE = regexp.MustCompile(`"method":\s*"([^"]+)"`)
	findCallSiteRE   = regexp.MustCompile(`"callSite":\s*"([^"]+)"`)
)

func IsInternalMethod(method string) bool {
	switch method {
	case "funcTestContract.upload", "funcTestContract.callConstructor", "funcTestContract.callMethod":
		return true
	case "contract.registerNode", "contract.getNodeRef", "cert.get":
		return true
	default:
		return false
	}
}

func (rv *RequestValidator) Validate(ctx context.Context, httpReq *http.Request) error {
	// all json rpc are POSTs, it's handled in other place
	if strings.ToUpper(httpReq.Method) != "POST" {
		return nil
	}

	// not RPC request, can not validate ATM
	if httpReq.URL.Path != "/admin-api/rpc" && httpReq.URL.Path != "/api/rpc" {
		return nil
	}

	logger := inslogger.FromContext(ctx)

	reqURL, err := url.Parse(httpReq.URL.String())
	if err != nil {
		logger.Error(errors.Wrap(err, "couldn't clone URL"))
		return errors.New("invalid URL")
	}

	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		return errors.New("failed to read request body")
	}

	httpReq.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	match := findCallSiteRE.FindSubmatch(body)
	if match == nil {
		match = findMethodNameRE.FindSubmatch(body)
		if match == nil {
			return errors.New("no 'method', 'callSite' or invalid JSON")
		}
	}

	if IsInternalMethod(string(match[1])) {
		return nil
	}

	reqURL.Path += "#" + string(match[1])

	route, pathParams, err := rv.router.FindRoute(httpReq.Method, reqURL)
	if err != nil {
		return errors.New("unknown method")
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	err = openapi3filter.ValidateRequest(ctx, requestValidationInput)
	if err != nil {
		return errors.Wrap(err, "request don't pass OpenAPI schema validation")
	}

	return nil
}
