package tlv

import (
	"bytes"
	"io"
	"reflect"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/packets"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) encoding.Encoder {
	return &Encoder{
		w: w,
	}
}

func (w *Encoder) Encode(p interface{}) error {
	var t TLV
	var err error
	switch p := p.(type) {
	case Marshaler:
		t, err = p.MarshalTLV()
	case *packets.Interest:
		t, err = marshalInterestPacket(p)
	case *packets.Data:
		t, err = marshalDataPacket(p)
	case TLV:
		t = p
	default:
		err = &encoding.UnsupportedTypeError{Encoding: "tlv", Type: reflect.TypeOf(p)}
	}
	if err != nil {
		return err
	}

	// blagh encoding this then copying is inefficient, but I don't know how to
	// check the maximum size otherwise.
	buf := &bytes.Buffer{}
	_, err = t.WriteTo(buf)
	if err != nil {
		return err
	}

	if len(buf.Bytes()) > encoding.MaxNDNPacketSize {
		return &encoding.PacketTooLargeError{}
	}

	_, err = io.Copy(w.w, buf)
	return err
}
