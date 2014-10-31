package ndn

import (
	"bytes"
	"errors"
	"io"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/encoding/tlv"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
)

var ErrUnexpectedEventTypeReceived = errors.New("unexpected event type received")
var ErrInvalidPacket = errors.New("invalid packet received")

func NewFace(transport Transport) *Face {
	return &Face{
		transport:            transport,
		encoder:              tlv.NewEncoder(transport),
		pendingInterestTable: newPendingInterestTable(),
	}
}

type Face struct {
	transport                Transport
	encoder                  encoding.Encoder
	pendingInterestTable     *pendingInterestTable
	commandKeyChain          *KeyChain
	commandCertificateName   name.Name
	commandInterestGenerator commandInterestGenerator
}

func (f *Face) SetCommandSigningInfo(keyChain *KeyChain, certificateName name.Name) {
	f.commandKeyChain = keyChain
	f.commandCertificateName = certificateName.Copy()
}

func (f *Face) ExpressInterest(i *packets.Interest) (*pendingInterest, error) {
	pendingInterest := f.pendingInterestTable.AddInterest(i)
	err := f.encoder.Encode(i)
	if err != nil {
		return nil, err
	}
	return pendingInterest, nil
}

func (f *Face) RemovePendingInterest(id uint64) {
	f.pendingInterestTable.RemovePendingInterest(id)
}

func (f *Face) ProcessEvents() error {
	r := tlv.NewDecoder(f.transport)

	for {
		var packet interface{}
		err := r.Decode(&packet)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		err = f.dispatchEvent(packet)
		if err != nil {
			return err
		}
	}
}

func (f *Face) dispatchEvent(packet interface{}) error {
	switch packet := packet.(type) {
	case *packets.Interest:
		// TODO: implement this
	case *packets.Data:
		// parent, ok := event.(tlv.ParentTLV)
		// if !ok /*|| len(parent.V) != 4 */ {
		// 	return ErrInvalidPacket
		// }
		// name, err := name.FromTLV(parent.V[0])
		// if err != nil {
		// 	return err
		// }
		f.pendingInterestTable.DispatchData(packet)
	default:
		return ErrUnexpectedEventTypeReceived
	}
	return nil
}

func (f *Face) Put(d *packets.Data) error {
	return f.encoder.Encode(d)
}

func (f *Face) RegisterPrefix(prefix name.Name) (<-chan *packets.Interest, error) {
	if f.commandKeyChain == nil {
		return nil, ErrCommandKeyChainNotSet
	}
	if f.commandCertificateName.Size() == 0 {
		return nil, ErrCommandCertificateNameNotSet
	}

	n := name.New(name.Component{"localhost"}, name.Component{"nfd"}, name.Component{"rib"}, name.Component{"register"})
	n = n.AppendBytes(f.registerPrefixControlParameters(prefix))
	commandInterest := &packets.Interest{}
	commandInterest.SetName(n)
	f.makeCommandInterest(commandInterest)
	// The interest is answered by the local host, so set a short timeout.
	commandInterest.SetInterestLifetime(2 * time.Second)
	// make command interest
	return nil, nil
}

func (f *Face) makeCommandInterest(i *packets.Interest) {
	f.commandInterestGenerator.generate(i, f.commandKeyChain, f.commandCertificateName, tlv.NewEncoder)
}

func (f *Face) registerPrefixControlParameters(prefix name.Name) []byte {
	buf := &bytes.Buffer{}
	(&controlParameters{name: prefix}).WriteTo(buf)
	return buf.Bytes()
}

func (f *Face) Close() error {
	return f.transport.Close()
}
