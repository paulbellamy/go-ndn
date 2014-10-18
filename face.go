package ndn

import (
	"bytes"
	"errors"
	"io"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
)

var ErrUnexpectedEventTypeReceived = errors.New("unexpected event type received")
var ErrInvalidPacket = errors.New("invalid packet received")

func NewFace(transport Transport) *Face {
	return &Face{
		transport:            transport,
		pendingInterestTable: newPendingInterestTable(),
	}
}

type Face struct {
	transport                Transport
	pendingInterestTable     *pendingInterestTable
	commandKeyChain          *KeyChain
	commandCertificateName   Name
	commandInterestGenerator commandInterestGenerator
}

func (f *Face) SetCommandSigningInfo(keyChain *KeyChain, certificateName Name) {
	f.commandKeyChain = keyChain
	f.commandCertificateName = certificateName.Copy()
}

func (f *Face) ExpressInterest(i *Interest) (*pendingInterest, error) {
	pendingInterest := f.pendingInterestTable.AddInterest(i)
	_, err := i.WriteTo(f.transport)
	if err != nil {
		return nil, err
	}
	return pendingInterest, nil
}

func (f *Face) RemovePendingInterest(id uint64) {
	f.pendingInterestTable.RemovePendingInterest(id)
}

func (f *Face) ProcessEvents() error {
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

func (f *Face) dispatchEvent(event encoding.TLV) error {
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

func (f *Face) Put(d *Data) error {
	// blagh encoding this then copying is inefficient, but I don't know how to
	// check the maximum size otherwise.
	buf := &bytes.Buffer{}
	_, err := d.WriteTo(buf)
	if err != nil {
		return err
	}

	if len(buf.Bytes()) > MaxNDNPacketSize {
		return PacketTooLargeError
	}

	_, err = io.Copy(f.transport, buf)
	return err
}

func (f *Face) RegisterPrefix(prefix Name) (<-chan *Interest, error) {
	if f.commandKeyChain == nil {
		return nil, ErrCommandKeyChainNotSet
	}
	if f.commandCertificateName.Size() == 0 {
		return nil, ErrCommandCertificateNameNotSet
	}

	name := Name{Component{"localhost"}, Component{"nfd"}, Component{"rib"}, Component{"register"}}
	name = name.AppendBytes(f.registerPrefixControlParameters(prefix))
	commandInterest := &Interest{name: name}
	f.makeCommandInterest(commandInterest)
	// The interest is answered by the local host, so set a short timeout.
	commandInterest.SetInterestLifetime(2 * time.Second)
	// make command interest
	return nil, nil
}

func (f *Face) makeCommandInterest(i *Interest) {
	f.commandInterestGenerator.generate(i, f.commandKeyChain, f.commandCertificateName)
}

func (f *Face) registerPrefixControlParameters(prefix Name) []byte {
	buf := &bytes.Buffer{}
	(&controlParameters{name: prefix}).WriteTo(buf)
	return buf.Bytes()
}

func (f *Face) Close() error {
	return f.transport.Close()
}
