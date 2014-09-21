package ndn

import (
	"errors"
	"io"

	"github.com/paulbellamy/go-ndn/encoding"
)

var ErrUnexpectedEventTypeReceived = errors.New("unexpected event type received")
var ErrInvalidPacket = errors.New("invalid packet received")

func Face(transport Transport) *face {
	return &face{
		transport:            transport,
		pendingInterestTable: newPendingInterestTable(),
	}
}

type face struct {
	transport            Transport
	pendingInterestTable *pendingInterestTable
}

func (f *face) ExpressInterest(i *Interest) (*pendingInterest, error) {
	pendingInterest := f.pendingInterestTable.AddInterest(i)
	_, err := i.WriteTo(f.transport)
	if err != nil {
		return nil, err
	}
	return pendingInterest, nil
}

func (f *face) RemovePendingInterest(id uint64) {
	f.pendingInterestTable.RemovePendingInterest(id)
}

func (f *face) ProcessEvents() error {
	r := encoding.NewReader(f.transport)

	for {
		tlv, err := r.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		err = f.dispatchEvent(tlv)
		if err != nil {
			return err
		}
	}
}

func (f *face) dispatchEvent(event encoding.TLV) error {
	switch event.Type() {
	case encoding.InterestType:
		// TODO: implement this
	case encoding.DataType:
		parent, ok := event.(encoding.ParentTLV)
		if !ok /*|| len(parent.V) != 4 */ {
			return ErrInvalidPacket
		}
		name, err := NameFromTLV(parent.V[0])
		if err != nil {
			return err
		}
		f.pendingInterestTable.DispatchData(&Data{name: name})
	default:
		return ErrUnexpectedEventTypeReceived
	}
	return nil
}

func (f *face) Close() error {
	return f.transport.Close()
}
