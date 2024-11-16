package WebsocketRouter

import (
	"encoding/json"
	"fmt"
	"strings"
)

func parsePath(path string) []string {
	return strings.SplitN(path, "/", 2)
}

func asyncEndpointRunner(endpoint func(payload []byte, companyID int) ([]byte, error), payload []byte, companyID int, responseChan chan []byte, errorChan chan error) {
	response, err := endpoint(payload, companyID)
	responseChan <- response
	errorChan <- err
}

func (router *RouteRegistration) HandleRequest(path string, payload []byte) ([]byte, error) {

	fmt.Println("Running router")

	parsedPath := parsePath(path)

	if len(parsedPath) == 1 {
		endpoint, ok := router.EndPoints[parsedPath[0]]
		if ok {
			fmt.Println("hi")

			responseChan := make(chan []byte)
			errChan := make(chan error)

			go asyncEndpointRunner(endpoint, payload, 123, responseChan, errChan)

			response := <-responseChan
			err := <-errChan

			return response, err

		}
	} else {
		subRouter, ok := router.SubRouters[parsedPath[0]]
		if ok {
			jsonBytes, _ := json.Marshal("abcd")
			return subRouter.HandleRequest(parsedPath[1], jsonBytes)
		}

	}

	//fmt.Printf("%+v\n", parsedPath[0])
	//fmt.Printf("%+v\n", parsedPath[1])
	return nil, nil
}
