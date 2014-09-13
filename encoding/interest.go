package encoding

import (
	"errors"
	"io"
)

var ErrUnexpectedTLVType = errors.New("Unexpected TLV Type")

type Interest struct {
	Name  string
	Nonce uint64
}

func (i *Interest) ReadFrom(r io.Reader) (int64, error) {
	n0, err := i.readName(r)
	if err != nil {
		return n0, err
	}
	n1, err := i.readNonce(r)
	return n0 + n1, err
}

func (i *Interest) readName(r io.Reader) (n int64, err error) {
	tlv := &TLV{}
	n, err = tlv.ReadFrom(r)
	if err != nil {
		return
	}

	if tlv.Type != NameType {
		err = ErrUnexpectedTLVType
		return
	}

	i.Name = string(tlv.Value)
	return
}

func (i *Interest) readNonce(r io.Reader) (n int64, err error) {
	tlv := &UintTLV{}
	n, err = tlv.ReadFrom(r)
	if err != nil {
		return
	}

	if tlv.Type != NonceType {
		err = ErrUnexpectedTLVType
		return
	}

	i.Nonce = tlv.Value
	return
}
