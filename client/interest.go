package client

import (
	"errors"
	"io"

	"github.com/paulbellamy/go-ndn/encoding"
)

var ErrInterestNameRequired = errors.New("interest name is required")

type Interest struct {
	Name Name
}

func (i *Interest) WriteTo(w io.Writer) (int64, error) {
	if i.Name.IsBlank() {
		return 0, ErrInterestNameRequired
	}

	return i.toTLV().WriteTo(w)
}

func (i *Interest) toTLV() encoding.TLV {
	return encoding.ParentTLV(encoding.InterestType, i.Name.toTLV())
}
