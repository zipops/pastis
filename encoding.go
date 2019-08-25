package pastis

import (
	"encoding/json"
	"io"
)

type Encoder interface {
	Encode(io.Writer, interface{}) error
}
type Decoder interface {
	Decode(io.Reader, interface{}) error
}

type EncodingJSON struct{}

func (EncodingJSON) Decode(r io.Reader, d interface{}) error {
	return json.NewDecoder(r).Decode(d)
}

func (EncodingJSON) Encode(r io.Writer, d interface{}) error {
	return json.NewEncoder(r).Encode(d)
}
