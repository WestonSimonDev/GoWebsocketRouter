package WebsocketRouter

var registeredRoutes = make(map[string]RouteRegistration)

type RouteRegistration struct {
	SubRouters map[string]*RouteRegistration
	EndPoints  map[string]func(payload []byte, companyID int) ([]byte, error)
}

func (rr *RouteRegistration) CreateSubRouter(pathPrefix string) (*RouteRegistration, error) {

	var newSubRouter = &RouteRegistration{make(map[string]*RouteRegistration), make(map[string]func(payload []byte, companyID int) ([]byte, error))}

	rr.SubRouters[pathPrefix] = newSubRouter

	return newSubRouter, nil
}

func (rr *RouteRegistration) CreateEndpoint(path string, endpoint func(payload []byte, companyID int) ([]byte, error)) {
	rr.EndPoints[path] = endpoint
}

func CreateToplevelRouter() (*RouteRegistration, error) {

	var newRouter = &RouteRegistration{make(map[string]*RouteRegistration), make(map[string]func(payload []byte, companyID int) ([]byte, error))}

	return newRouter, nil

}
