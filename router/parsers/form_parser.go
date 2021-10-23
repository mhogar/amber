package parsers

import (
	"net/http"

	"github.com/mhogar/amber/common"
	"github.com/monoculum/formam/v3"
)

type FormBodyParser struct{}

func NewFormBodyParser() FormBodyParser {
	return FormBodyParser{}
}

func (FormBodyParser) ParseBody(req *http.Request, v interface{}) error {
	//parse the form
	err := req.ParseForm()
	if err != nil {
		return common.ChainError("error parsing form", err)
	}

	//decode into interface
	decoder := formam.NewDecoder(&formam.DecoderOptions{TagName: "json"})
	err = decoder.Decode(req.Form, v)
	if err != nil {
		return common.ChainError("error decoding form data", err)
	}

	return nil
}
