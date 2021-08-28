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
	KeyLoader   loaders.RSAKeyLoader
	TokenSigner TokenSigner
}

func (tfs CoreTokenFactorySelector) Select(tokenType int) TokenFactory {
	if tokenType == TokenTypeFirebase {
		return &FirebaseTokenFactory{
			JSONLoader:  tfs.JSONLoader,
			KeyLoader:   tfs.KeyLoader,
			TokenSigner: tfs.TokenSigner,
		}
	}

	return nil
}
