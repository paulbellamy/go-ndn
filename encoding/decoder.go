package encoding

import "io"

type Decoder interface {
	Decode() (interface{}, error)
}

type DecoderFactory func(io.Reader) Decoder
