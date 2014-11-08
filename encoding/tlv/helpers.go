package tlv

import (
	"bytes"
	"fmt"
	"io"

	"github.com/paulbellamy/go-ndn/encoding"
)

type parser interface {
	Parse([]byte) (interface{}, []byte, error)
}

type parserFunc func([]byte) (interface{}, []byte, error)

func (f parserFunc) Parse(input []byte) (interface{}, []byte, error) {
	return f(input)
}

var nonNegativeInteger = parserFunc(func(input []byte) (interface{}, []byte, error) {
	value, numRead, err := ReadUint(bytes.NewReader(input))
	if err != nil {
		return nil, input, io.ErrUnexpectedEOF
	}
	return value, input[numRead:], nil
})

var Byte = parserFunc(func(input []byte) (interface{}, []byte, error) {
	if len(input) <= 0 {
		return nil, input, io.ErrUnexpectedEOF
	}
	return input[0], input[1:], nil
})

var empty = parserFunc(func(input []byte) (interface{}, []byte, error) {
	if len(input) > 0 {
		return nil, input, &encoding.InvalidUnmarshalError{Message: "unexpected bytes"}
	}
	return nil, input, nil
})

type bytesParserType func(int) parser

var Bytes bytesParserType = func(n int) parser {
	return bytesWithLengthParser(n)
}

func (b bytesParserType) Parse(input []byte) (interface{}, []byte, error) {
	return input, []byte{}, nil
}

type bytesWithLengthParser int

func (b bytesWithLengthParser) Parse(input []byte) (interface{}, []byte, error) {
	if len(input) < int(b) {
		return nil, input, io.ErrUnexpectedEOF
	}
	return input[:b], input[b:], nil
}

type seqParser []parser

// A sequence of things which must all match
func seq(ps ...parser) parser {
	if len(ps) == 1 {
		return ps[0]
	}
	return seqParser(ps)
}

func (s seqParser) Parse(input []byte) (interface{}, []byte, error) {
	var result = []interface{}{}
	var rest = input
	var err error
	for _, p := range s {
		var item interface{}
		item, rest, err = p.Parse(rest)
		if err != nil {
			return nil, input, err
		}

		if item != nil {
			result = append(result, item)
		}
	}
	return result, rest, nil
}

func times(n int, p parser) parser {
	ps := []parser{}
	for i := 0; i < n; i++ {
		ps = append(ps, p)
	}
	return seq(ps...)
}

func exactly(bs ...byte) parser {
	return exactlyParser(bs)
}

type exactlyParser []byte

func (p exactlyParser) Parse(input []byte) (interface{}, []byte, error) {
	if len(input) < len(p) {
		return nil, input, io.ErrUnexpectedEOF
	}
	for i, b := range p {
		if input[i] != b {
			return nil, input, io.ErrUnexpectedEOF
		}
	}
	return input[:len(p)], input[len(p):], nil
}

type tlvParser struct {
	t uint64
	p parser
}

func tlv(t uint64, ps ...parser) parser {
	return tlvParser{t: t, p: seq(ps...)}
}

func (p tlvParser) Parse(input []byte) (interface{}, []byte, error) {
	r := bytes.NewReader(input)

	t, n, err := ReadNumber(r)
	if err != nil {
		return nil, input, err
	}

	if t != p.t {
		return nil, input, &encoding.InvalidUnmarshalError{Message: fmt.Sprintf("unexpected tlv type %d, expected %d", t, p.t)}
	}

	length, read, err := ReadNumber(r)
	n += read
	if err != nil {
		return nil, input, err
	}

	end := int64(length) + n
	if end > int64(len(input)) {
		return nil, input, io.ErrUnexpectedEOF
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
	return zeroOrMoreParser{seq(ps...)}
}

type zeroOrMoreParser struct {
	parser
}

func (p zeroOrMoreParser) Parse(input []byte) (interface{}, []byte, error) {
	found := []interface{}{}
	rest := input
	var err error
	for {
		var item interface{}
		item, rest, err = p.parser.Parse(rest)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				break
			}
			return nil, input, err
		}
		found = append(found, item)
	}
	return found, rest, nil
}

// one-or-more repetitions of this seq
func oneOrMore(ps ...parser) parser {
	return oneOrMoreParser{seq(ps...)}
}

type oneOrMoreParser struct {
	parser
}

func (p oneOrMoreParser) Parse(input []byte) (interface{}, []byte, error) {
	found := []interface{}{}
	rest := input
	var err error
	for {
		var item interface{}
		item, rest, err = p.parser.Parse(rest)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				break
			}
			return nil, input, err
		}
		found = append(found, item)
	}
	if len(found) == 0 {
		return nil, input, io.ErrUnexpectedEOF
	}
	return found, rest, nil
}

// first parser to match
func or(p parser, ps ...parser) parser {
	result := orParser{p}
	result = append(result, ps...)
	return result
}

type orParser []parser

func (ps orParser) Parse(input []byte) (interface{}, []byte, error) {
	var item interface{}
	var rest []byte
	var err error
	for _, p := range ps {
		item, rest, err = p.Parse(input)
		if err == nil {
			return item, rest, err
		}
	}
	return nil, input, io.ErrUnexpectedEOF
}

// zero-or-one matches
func maybe(p parser) parser {
	return parserFunc(func(input []byte) (interface{}, []byte, error) {
		result, rest, err := p.Parse(input)
		if err != nil {
			return nil, input, nil
		}
		return result, rest, err
	})
}
