package models

import (
	"net/url"

	"github.com/google/uuid"
)

const (
	ValidateClientValid              = 0x0
	ValidateClientNilUID             = 0x1
	ValidateClientEmptyName          = 0x2
	ValidateClientNameTooLong        = 0x4
	ValidateClientEmptyRedirectUrl   = 0x8
	ValidateClientRedirectUrlTooLong = 0x10
	ValidateClientInvalidRedirectUrl = 0x20
)

// ClientNameMaxLength is the max length a client's name can be.
const ClientNameMaxLength = 30

// ClientRedirectUrlMaxLength is the max length a client's redirect url can be.
const ClientRedirectUrlMaxLength = 100

// Client represents the client model.
type Client struct {
	UID         uuid.UUID
	Name        string
	RedirectUrl string
}

type ClientCRUD interface {
	// CreateClient creates a new client and returns any errors.
	CreateClient(client *Client) error

	// GetClientByUID fetches the client associated with the uid.
	// If no clients are found, returns nil client. Also returns any errors.
	GetClientByUID(uid uuid.UUID) (*Client, error)

	// UpdateClient updates the client.
	// Returns result of whether the client was found, and any errors.
	UpdateClient(client *Client) (bool, error)

	// DeleteClient deletes the client the with the given uid.
	// Returns result of whether the client was found, and any errors.
	DeleteClient(uid uuid.UUID) (bool, error)
}

func CreateClient(uid uuid.UUID, name string, redirectUrl string) *Client {
	return &Client{
		UID:         uid,
		Name:        name,
		RedirectUrl: redirectUrl,
	}
}

func CreateNewClient(name string, redirectUrl string) *Client {
	return CreateClient(uuid.New(), name, redirectUrl)
}

// Validate validates the client model has valid fields.
// Returns an int indicating which fields are invalid.
func (c *Client) Validate() int {
	code := ValidateClientValid

	if c.UID == uuid.Nil {
		code |= ValidateClientNilUID
	}

	if c.Name == "" {
		code |= ValidateClientEmptyName
	} else if len(c.Name) > ClientNameMaxLength {
		code |= ValidateClientNameTooLong
	}

	if c.RedirectUrl == "" {
		code |= ValidateClientEmptyRedirectUrl
	} else if len(c.RedirectUrl) > ClientRedirectUrlMaxLength {
		code |= ValidateClientRedirectUrlTooLong
	} else {
		_, err := url.Parse(c.RedirectUrl)
		if err != nil {
			code |= ValidateClientInvalidRedirectUrl
		}
	}

	return code
}
