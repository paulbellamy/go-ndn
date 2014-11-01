package tlv

import (
	"bytes"
	"io"

	"github.com/paulbellamy/go-ndn/encoding"
)

type parser interface {
	Parse([]byte) (interface{}, []byte, error)
}

type parserFunc func([]byte) (interface{}, []byte, error)

func (f parserFunc) Parse(body []byte) (interface{}, []byte, error) {
	return f(body)
}

var nonNegativeInteger = parserFunc(func(body []byte) (interface{}, []byte, error) {
	value, numRead, err := ReadUint(bytes.NewReader(body))
	if err != nil {
		return nil, body, io.ErrUnexpectedEOF
	}
	return value, body[numRead:], nil
})

var Byte = parserFunc(func(body []byte) (interface{}, []byte, error) {
	if len(body) <= 0 {
		return nil, body, io.ErrUnexpectedEOF
	}
	return body[0], body[1:], nil
})

var empty = parserFunc(func(body []byte) (interface{}, []byte, error) {
	if len(body) > 0 {
		return nil, body, &encoding.InvalidUnmarshalError{Message: "unexpected bytes"}
	}
	return nil, body, nil
})

type bytesParserType func(int) parser

var Bytes bytesParserType = func(n int) parser {
	return bytesWithLengthParser(n)
}

func (b bytesParserType) Parse(body []byte) (interface{}, []byte, error) {
	return body, []byte{}, nil
}

type bytesWithLengthParser int

func (b bytesWithLengthParser) Parse(body []byte) (interface{}, []byte, error) {
	if len(body) < int(b) {
		return nil, body, io.ErrUnexpectedEOF
	}
	return body[:b], body[b:], nil
}

type seqParser []parser

// A sequence of things which must all match
func seq(ps ...parser) parser {
	return seqParser(ps)
}

func (s seqParser) Parse(body []byte) (interface{}, []byte, error) {
	var result = []interface{}{}
	var rest = body
	var err error
	for _, p := range s {
		var item interface{}
		item, rest, err = p.Parse(rest)
		if err != nil {
			return nil, body, err
		}

		result = append(result, item)
	}
	return result, rest, nil
}

func times(n int, p parser) parser {
	return seq()
}

func exactly(bs ...byte) parser {
	return seq()
}

type tlvParser struct {
	t uint64
	p parser
}

func tlv(t uint64, ps ...parser) parser {
	var p parser
	if len(ps) > 1 {
		p = seq(ps...)
	} else if len(ps) == 1 {
		p = ps[0]
	}
	return tlvParser{t: t, p: p}
}

func (p tlvParser) Parse(input []byte) (interface{}, []byte, error) {
	r := bytes.NewReader(input)

	t, n, err := ReadNumber(r)
	if err != nil {
		return nil, input, err
	}

	length, read, err := ReadNumber(r)
	n += read
	if err != nil {
		return nil, input, err
	}

	end := int64(len(input[n:]))
	if n+int64(length) < end {
		end = n + int64(length)
	}
	value, rest, err := p.p.Parse(input[n:end])
	if err != nil {
		return nil, input, err
	}

	if len(rest) != 0 {
		return nil, input, &encoding.InvalidUnmarshalError{Message: "unexpected trailing bytes"}
	}

	return GenericTLV{T: t, V: value}, input[end:], nil
}

// zero-or-more repetitions of this seq
func zeroOrMore(ps ...parser) parser {
	return seq(ps...)
}

// one-or-more repetitions of this seq
func oneOrMore(ps ...parser) parser {
	return seq(ps...)
}

// first parser to match
func or(p parser, ps ...parser) parser {
	return p
}

// zero-or-one matches
func maybe(p parser) parser {
	return p
}
