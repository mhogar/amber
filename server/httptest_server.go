package server

import (
	"authserver/router"
	"net/http/httptest"
)

// HTTPTestServer is a wrapper for an httptest server that implements the server interface.
type HTTPTestServer struct {
	*httptest.Server
}

func CreateHTTPTestServerRunner(routerFactory router.RouterFactory) Runner {
	return Runner{
		Server: &HTTPTestServer{
			Server: httptest.NewUnstartedServer(routerFactory.CreateRouter()),
		},
	}
}

// Start start the server. Always returns a nil error.
func (s *HTTPTestServer) Start() error {
	s.Server.Start()
	return nil
}
