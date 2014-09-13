package client

import (
	"errors"

	"github.com/paulbellamy/go-ndn/encoding"
)

var ErrUnexpectedEventTypeReceived = errors.New("unexpected event type received")

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
		event, err := r.Read()
		if err != nil {
			return err
		}

		return f.dispatchEvent(event)
	}
}

func (f *face) dispatchEvent(event interface{}) error {
	switch event := event.(type) {
	case *Interest:
	case *Data:
		f.pendingInterestTable.DispatchData(event)
	default:
		return ErrUnexpectedEventTypeReceived
	}
	return nil
}

func (f *face) Close() error {
	return f.transport.Close()
}
