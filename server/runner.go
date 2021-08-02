package server

type Server interface {
	// Start starts the server and returns any errors encountered while it is running
	Start() error

	// Close closes the server
	Close()
}

type Runner struct {
	Server Server
}

// Run runs the server and returns any errors
func (s Runner) Run() error {
	return s.Server.Start()
}
