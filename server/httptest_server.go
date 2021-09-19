package server

import (
	"net/http/httptest"

	"github.com/mhogar/amber/router"
)

// HTTPTestServer is a wrapper for an httptest server that implements the server interface.
type HTTPTestServer struct {
	*httptest.Server
}

// CreateHTTPTestServerRunner creates a new Runner using an HTTPTestServer.
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
