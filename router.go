package WebsocketRouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CreateToplevelRouter() (*TopLevelRouteRegistration, error) {

	var newRouter = &TopLevelRouteRegistration{make(map[string]*RouteRegistration), make(map[string]func(payload []byte, httpRequest *http.Request) ([]byte, error))}

	newRouter.CreateEndpoint("heartBeat", func(payload []byte, httpRequest *http.Request) ([]byte, error) {
		var pathNotFoundResponse = make(map[string]string)
		pathNotFoundResponse["error"] = "successful"
		response, _ := json.Marshal(pathNotFoundResponse)
		return response, nil
	})

	newRouter.CreateSubRouter("response")

	return newRouter, nil

}

func parsePath(path string) []string {
	return strings.SplitN(path, "/", 2)
}

func asyncEndpointRunner(endpoint func(payload []byte, httpRequest *http.Request) ([]byte, error), payload []byte, responseChan chan []byte, errorChan chan error, httpRequest *http.Request) {
	response, err := endpoint(payload, httpRequest)
	responseChan <- response
	errorChan <- err
}

func (router *RouteRegistration) HandleRequest(request Request, httpRequest *http.Request) ([]byte, error) {

	parsedPath := parsePath(request.Path)

	request.Path = parsedPath[len(parsedPath)-1]

	if len(parsedPath) == 1 {
		endpoint, ok := router.EndPoints[parsedPath[0]]
		if ok {

			responseChan := make(chan []byte)
			errChan := make(chan error)

			go asyncEndpointRunner(endpoint, request.Payload, responseChan, errChan, httpRequest)

			response := <-responseChan
			err := <-errChan
			responseBytes, jsonErr := json.Marshal(Response{Path: fmt.Sprintf("response/%s", request.Response), Payload: response})
			if jsonErr != nil {
				return nil, jsonErr
			}
			return responseBytes, err

		} else {
			wildCardEndpoint, wildOk := router.EndPoints["*"]
			if wildOk {
				responseChan := make(chan []byte)
				errChan := make(chan error)

				go asyncEndpointRunner(wildCardEndpoint, request.Payload, responseChan, errChan, httpRequest)

				response := <-responseChan
				err := <-errChan
				responseBytes, jsonErr := json.Marshal(Response{Path: fmt.Sprintf("response/%s", request.Response), Payload: response})
				if jsonErr != nil {
					return nil, jsonErr
				}
				return responseBytes, err
			}
		}
	} else {
		subRouter, ok := router.SubRouters[parsedPath[0]]
		if ok {
			return subRouter.HandleRequest(request, httpRequest)
		}

	}
	var pathNotFoundResponse = make(map[string]string)
	pathNotFoundResponse["error"] = "Action Not Found"
	payload, _ := json.Marshal(pathNotFoundResponse)

	responseBytes, jsonErr := json.Marshal(Response{Path: fmt.Sprintf("response/%s", request.Response), Payload: payload})
	if jsonErr != nil {
		return nil, jsonErr
	}
	return responseBytes, nil
}
