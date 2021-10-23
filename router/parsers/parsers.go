package parsers

import "net/http"

type BodyParser interface {
	ParseBody(req *http.Request, v interface{}) error
}
