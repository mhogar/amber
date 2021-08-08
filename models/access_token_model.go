package models

import "github.com/google/uuid"

const (
	ValidateAccessTokenValid         = 0x0
	ValidateAccessTokenNilID         = 0x1
	ValidateAccessTokenNilUser       = 0x2
	ValidateAccessTokenInvalidUser   = 0x4
	ValidateAccessTokenNilClient     = 0x8
	ValidateAccessTokenInvalidClient = 0x10
)

// AccessToken represents the access token model
type AccessToken struct {
	ID     uuid.UUID
	User   *User
	Client *Client
}

type AccessTokenCRUD interface {
	// SaveAccessToken saves the access token and returns any errors
	SaveAccessToken(token *AccessToken) error

	// GetAccessTokenByID fetches the access token associated with the id
	// If no tokens are found, returns nil token. Also returns any errors
	GetAccessTokenByID(ID uuid.UUID) (*AccessToken, error)

	// DeleteAccessToken deletes the token and returns any errors
	DeleteAccessToken(token *AccessToken) error

	// DeleteAllOtherUserTokens deletes all of the user's tokens expect for the provided one and returns any errors
	DeleteAllOtherUserTokens(token *AccessToken) error
}

func CreateAccessToken(id uuid.UUID, user *User, client *Client) *AccessToken {
	return &AccessToken{
		ID:     id,
		User:   user,
		Client: client,
	}
}

func CreateNewAccessToken(user *User, client *Client) *AccessToken {
	return CreateAccessToken(uuid.New(), user, client)
}

// Validate validates the access token model has valid fields
// Returns an int indicating which fields are invalid
func (tk *AccessToken) Validate() int {
	code := ValidateAccessTokenValid

	if tk.ID == uuid.Nil {
		code |= ValidateAccessTokenNilID
	}

	if tk.User == nil {
		code |= ValidateAccessTokenNilUser
	} else {
		verr := tk.User.Validate()
		if verr != ValidateUserValid {
			code |= ValidateAccessTokenInvalidUser
		}
	}

	if tk.Client == nil {
		code |= ValidateAccessTokenNilClient
	} else {
		verr := tk.Client.Validate()
		if verr != ValidateClientValid {
			code |= ValidateAccessTokenInvalidClient
		}
	}

	return code
}
