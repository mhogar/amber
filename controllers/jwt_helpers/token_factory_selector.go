package jwthelpers

import (
	"github.com/mhogar/amber/loaders"
	"github.com/mhogar/amber/models"
)

type TokenFactorySelector interface {
	// Select selects the TokenFactory based on the provided tokenType.
	// Returns the TokenFactory, or nil if the type is unknown.
	Select(tokenType int) TokenFactory
}

type CoreTokenFactorySelector struct {
	JSONLoader  loaders.JSONLoader
	DataLoader  loaders.RawDataLoader
	TokenSigner TokenSigner
}

func (tfs CoreTokenFactorySelector) Select(tokenType int) TokenFactory {
	//default token type
	if tokenType == models.ClientTokenTypeDefault {
		return &DefaultTokenFactory{
			DataLoader:  tfs.DataLoader,
			TokenSigner: tfs.TokenSigner,
		}
	}

	//firebase token type
	if tokenType == models.ClientTokenTypeFirebase {
		return &FirebaseTokenFactory{
			JSONLoader:  tfs.JSONLoader,
			TokenSigner: tfs.TokenSigner,
		}
	}

	//unknown token type
	return nil
}
