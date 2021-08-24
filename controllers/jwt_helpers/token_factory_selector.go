package jwthelpers

const (
	TokenTypeFirebase = iota
)

type TokenFactorySelector interface {
	// Select selects the TokenFactory based on the provided tokenType.
	// Returns the TokenFactory, or nilif the type is unknown.
	Select(tokenType int) TokenFactory
}

type CoreTokenFactorySelector struct{}

func (tfs CoreTokenFactorySelector) Select(tokenType int) TokenFactory {
	if tokenType == TokenTypeFirebase {
		return &FirebaseTokenFactory{}
	}

	return nil
}
