package models

import "github.com/google/uuid"

// AccessToken represents the access token model.
type AccessToken struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	ClientID uuid.UUID
	ScopeID  uuid.UUID
}

// AccessTokenCRUD is an interface for performing CRUD operations on an access token.
type AccessTokenCRUD interface {
	// SaveAccessToken saves the access token and returns any errors.
	SaveAccessToken(token *AccessToken) error

	// GetAccessTokenByID fetches the access token associated with the id.
	// If no tokens are found, returns nil token. Also returns any errors.
	GetAccessTokenByID(ID uuid.UUID) (*AccessToken, error)

	// DeleteAccessToken deletes the token and returns any errors.
	DeleteAccessToken(token *AccessToken) error
}

// CreateNewAccessToken creates a access token model with a new id and the provided fields.
func CreateNewAccessToken(userID uuid.UUID, clientID uuid.UUID, scopeID uuid.UUID) *AccessToken {
	return &AccessToken{
		ID:       uuid.New(),
		UserID:   userID,
		ClientID: clientID,
		ScopeID:  scopeID,
	}
}
