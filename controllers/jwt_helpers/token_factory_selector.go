package jwthelpers

import "authserver/loaders"

const (
	TokenTypeDefault  = iota
	TokenTypeFirebase = iota
)

type TokenFactorySelector interface {
	// Select selects the TokenFactory based on the provided tokenType.
	// Returns the TokenFactory, or nilif the type is unknown.
	Select(tokenType int) TokenFactory
}

type CoreTokenFactorySelector struct {
	JSONLoader  loaders.JSONLoader
	DataLoader  loaders.RawDataLoader
	TokenSigner TokenSigner
}

func (tfs CoreTokenFactorySelector) Select(tokenType int) TokenFactory {
	//default token type
	if tokenType == TokenTypeDefault {
		return &DefaultTokenFactory{
			DataLoader:  tfs.DataLoader,
			TokenSigner: tfs.TokenSigner,
		}
	}

	//firebase token type
	if tokenType == TokenTypeFirebase {
		return &FirebaseTokenFactory{
			JSONLoader:  tfs.JSONLoader,
			TokenSigner: tfs.TokenSigner,
		}
	}

	//unknown token type
	return nil
}
