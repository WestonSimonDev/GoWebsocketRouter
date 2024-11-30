package WebsocketRouter

import (
	"encoding/json"
	"net/http"
)

type Router interface {
	CreateSubRouter(pathPrefix string) (*RouteRegistration, error)
	CreateEndpoint(path string, endpoint func(payload []byte, httpRequest *http.Request) ([]byte, error))
}

type RouteRegistration struct {
	SubRouters map[string]*RouteRegistration
	EndPoints  map[string]func(payload []byte, httpRequest *http.Request) ([]byte, error)
}

type TopLevelRouteRegistration struct {
	SubRouters map[string]*RouteRegistration
	EndPoints  map[string]func(payload []byte, httpRequest *http.Request) ([]byte, error)
}

func (rr *TopLevelRouteRegistration) CreateResponseEndpoint(responseID string, endpoint func(payload []byte, httpRequest *http.Request) ([]byte, error)) {
	rr.SubRouters["response"].EndPoints[responseID] = endpoint
}

func (rr *TopLevelRouteRegistration) ConsumeResponseEndpoint(responseID string) {
	delete(rr.SubRouters["response"].EndPoints, responseID)
}

func (rr *TopLevelRouteRegistration) CreateSubRouter(pathPrefix string) (*RouteRegistration, error) {

	var newSubRouter = &RouteRegistration{make(map[string]*RouteRegistration), make(map[string]func(payload []byte, httpRequest *http.Request) ([]byte, error))}

	rr.SubRouters[pathPrefix] = newSubRouter

	return newSubRouter, nil
}

func (rr *TopLevelRouteRegistration) CreateEndpoint(path string, endpoint func(payload []byte, httpRequest *http.Request) ([]byte, error)) {
	rr.EndPoints[path] = endpoint
}

func (rr *RouteRegistration) CreateSubRouter(pathPrefix string) (*RouteRegistration, error) {

	var newSubRouter = &RouteRegistration{make(map[string]*RouteRegistration), make(map[string]func(payload []byte, httpRequest *http.Request) ([]byte, error))}

	rr.SubRouters[pathPrefix] = newSubRouter

	return newSubRouter, nil
}

func (rr *RouteRegistration) CreateEndpoint(path string, endpoint func(payload []byte, httpRequest *http.Request) ([]byte, error)) {
	rr.EndPoints[path] = endpoint
}

type Request struct {
	Path     string          `json:"action"`
	Payload  json.RawMessage `json:"payload"`
	Response string          `json:"response"`
}

func NewRequest(requestBytes []byte) (Request, error) {
	var request Request
	err := json.Unmarshal(requestBytes, &request)
	return request, err
}

type Response struct {
	Path    string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}
