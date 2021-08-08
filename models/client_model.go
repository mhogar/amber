package models

import (
	"github.com/google/uuid"
)

const (
	ValidateClientValid       = 0x0
	ValidateClientNilID       = 0x1
	ValidateClientEmptyName   = 0x2
	ValidateClientNameTooLong = 0x4
)

// ClientNameMaxLength is the max length a client's name can be
const ClientNameMaxLength = 30

// Client represents the client model
type Client struct {
	ID   uuid.UUID
	Name string
}

type ClientCRUD interface {
	// SaveClient saves the client and returns any errors
	SaveClient(client *Client) error

	// UpdateClient updates the client and returns any errors
	UpdateClient(client *Client) error

	// GetClientByID fetches the client associated with the id
	// If no clients are found, returns nil client. Also returns any errors
	GetClientByID(ID uuid.UUID) (*Client, error)
}

func CreateClient(id uuid.UUID, name string) *Client {
	return &Client{
		ID:   id,
		Name: name,
	}
}

func CreateNewClient(name string) *Client {
	return CreateClient(uuid.New(), name)
}

// Validate validates the client model has valid fields
// Returns an int indicating which fields are invalid
func (c *Client) Validate() int {
	code := ValidateClientValid

	if c.ID == uuid.Nil {
		code |= ValidateClientNilID
	}

	if c.Name == "" {
		code |= ValidateClientEmptyName
	} else if len(c.Name) > ClientNameMaxLength {
		code |= ValidateClientNameTooLong
	}

	return code
}
