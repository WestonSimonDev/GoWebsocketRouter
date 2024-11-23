package WebsocketRouter

import "encoding/json"

type RouteRegistration struct {
	SubRouters map[string]*RouteRegistration
	EndPoints  map[string]func(payload []byte) ([]byte, error)
}

func (rr *RouteRegistration) CreateSubRouter(pathPrefix string) (*RouteRegistration, error) {

	var newSubRouter = &RouteRegistration{make(map[string]*RouteRegistration), make(map[string]func(payload []byte) ([]byte, error))}

	rr.SubRouters[pathPrefix] = newSubRouter

	return newSubRouter, nil
}

func (rr *RouteRegistration) CreateEndpoint(path string, endpoint func(payload []byte) ([]byte, error)) {
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
