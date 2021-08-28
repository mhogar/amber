package loaders

type RawDataLoader interface {
	// Load loads the raw bytes from the provided uri.
	// Returns the bytes and any errors.
	Load(uri string) ([]byte, error)
}
