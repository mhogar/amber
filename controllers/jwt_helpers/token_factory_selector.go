package jwthelpers

import "authserver/loaders"

const (
	TokenTypeFirebase = iota
)

type TokenFactorySelector interface {
	// Select selects the TokenFactory based on the provided tokenType.
	// Returns the TokenFactory, or nilif the type is unknown.
	Select(tokenType int) TokenFactory
}

type CoreTokenFactorySelector struct {
	JSONLoader  loaders.JSONLoader
	TokenSigner TokenSigner
}

func (tfs CoreTokenFactorySelector) Select(tokenType int) TokenFactory {
	if tokenType == TokenTypeFirebase {
		return &FirebaseTokenFactory{
			JSONLoader:  tfs.JSONLoader,
			TokenSigner: tfs.TokenSigner,
		}
	}

	return nil
}
