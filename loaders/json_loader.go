package loaders

type JSONLoader interface {
	// Load loads the json from the provided uri into v.
	// Returns any errors.
	Load(uri string, v interface{}) error
}
