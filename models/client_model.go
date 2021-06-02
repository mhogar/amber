package models

import (
	"github.com/google/uuid"
)

// Client ValidateError statuses.
const (
	ValidateClientValid = 0x0
	ValidateClientNilID = 0x1
)

// Client represents the client model.
type Client struct {
	ID uuid.UUID
}

// ClientCRUD is an interface for performing CRUD operations on a client.
type ClientCRUD interface {
	// SaveClient saves the client and returns any errors.
	SaveClient(client *Client) error

	// GetClientByID fetches the client associated with the id.
	// If no clients are found, returns nil client. Also returns any errors.
	GetClientByID(ID uuid.UUID) (*Client, error)
}

// CreateNewClient creates a client model with a new id and the provided fields.
func CreateNewClient() *Client {
	return &Client{
		ID: uuid.New(),
	}
}

// Validate validates the client model has valid fields.
// Returns an int indicating which fields are invalid.
func (c *Client) Validate() int {
	code := ValidateClientValid

	if c.ID == uuid.Nil {
		code |= ValidateClientNilID
	}

	return code
}
