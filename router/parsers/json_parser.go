package parsers

import (
	"encoding/json"
	"net/http"

	"github.com/mhogar/amber/common"
)

type JSONBodyParser struct{}

func NewJSONBodyParser() JSONBodyParser {
	return JSONBodyParser{}
}

func (JSONBodyParser) ParseBody(req *http.Request, v interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(v)
	if err != nil {
		return common.ChainError("error decoding json", err)
	}

	return nil
}
