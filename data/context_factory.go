package data

import (
	"context"
	"time"
)

type ContextFactory struct {
	Context context.Context
	Timeout int
}

// CreateStandardTimeoutContext creates a context with the configured timeout.
// It is a child of the configured context and can be canceled by that context's cancel function.
// Returns the created context and its cancel function.
func (cf ContextFactory) CreateStandardTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(cf.Context, time.Duration(cf.Timeout)*time.Millisecond)
}
