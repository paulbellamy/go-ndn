package encoding

import "io"

type Encoder interface {
	Encode(v interface{}) error
}

type EncoderFactory func(io.Writer) Encoder
