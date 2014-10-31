package encoding

import "io"

type Decoder interface {
	Decode(v interface{}) error
}

type DecoderFactory func(io.Reader) Decoder
