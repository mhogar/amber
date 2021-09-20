package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mhogar/amber/router"
)

// HTTPServer is a wrapper for an http server that implements the server interface.
type HTTPServer struct {
	http.Server
}

// CreateHTTPServerRunner creates a new Runner using an HTTPServer.
func CreateHTTPServerRunner(routerFactory router.RouterFactory) Runner {
	//get port from env variable (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Runner{
		Server: &HTTPServer{
			Server: http.Server{
				Addr:    ":" + port,
				Handler: routerFactory.CreateRouter(),
			},
		},
	}
}

// Start starts the http server. Always returns a non-nil error.
func (s *HTTPServer) Start() error {
	fmt.Println("Server is running on port", s.Addr)
	return s.ListenAndServe()
}

// Close does nothing but exists to satisfy the server interface.
func (*HTTPServer) Close() {}
