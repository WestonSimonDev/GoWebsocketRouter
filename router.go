package WebsocketRouter

import (
	"encoding/json"
	"fmt"
	"strings"
)

func CreateToplevelRouter() (*RouteRegistration, error) {

	var newRouter = &RouteRegistration{make(map[string]*RouteRegistration), make(map[string]func(payload []byte) ([]byte, error))}

	newRouter.CreateEndpoint("heartBeat", func(payload []byte) ([]byte, error) {
		return []byte(""), nil
	})

	return newRouter, nil

}

func parsePath(path string) []string {
	return strings.SplitN(path, "/", 2)
}

func asyncEndpointRunner(endpoint func(payload []byte) ([]byte, error), payload []byte, responseChan chan []byte, errorChan chan error) {
	response, err := endpoint(payload)
	responseChan <- response
	errorChan <- err
}

func (router *RouteRegistration) HandleRequest(request Request) ([]byte, error) {

	fmt.Println("Running router")

	parsedPath := parsePath(request.Path)

	request.Path = parsedPath[len(parsedPath)-1]

	if len(parsedPath) == 1 {
		endpoint, ok := router.EndPoints[parsedPath[0]]
		if ok {
			fmt.Println("hi")

			responseChan := make(chan []byte)
			errChan := make(chan error)

			go asyncEndpointRunner(endpoint, request.Payload, responseChan, errChan)

			response := <-responseChan
			err := <-errChan
			responseBytes, jsonErr := json.Marshal(Response{Path: fmt.Sprintf("response/%s", request.Response), Payload: response})
			if jsonErr != nil {
				return nil, jsonErr
			}
			return responseBytes, err

		}
	} else {
		subRouter, ok := router.SubRouters[parsedPath[0]]
		if ok {
			return subRouter.HandleRequest(request)
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
