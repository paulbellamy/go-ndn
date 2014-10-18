package ndn

import (
	"io"

	"github.com/paulbellamy/go-ndn/encoding"
)

type Data struct {
	name Name
	MetaInfo
	content   []byte
	signature Signature
}

func (d *Data) GetName() Name {
	if d.name == nil {
		d.name = Name{}
	}
	return d.name
}

func (d *Data) SetName(x Name) {
	d.name = x.Copy()
}

func (d *Data) GetContent() []byte {
	if d.content == nil {
		d.content = []byte{}
	}
	return d.content
}

func (d *Data) SetContent(x []byte) {
	d.content = x
}

func (d *Data) GetSignature() Signature {
	return d.signature
}

func (d *Data) SetSignature(s Signature) {
	d.signature = s
}

func (d *Data) WriteTo(w io.Writer) (int64, error) {
	if d.GetName().IsBlank() {
		return 0, ErrNameRequired
	}

	return d.toTLV().WriteTo(w)
}

func (d *Data) toTLV() encoding.TLV {
	return encoding.ParentTLV{encoding.DataType, []encoding.TLV{d.GetName().toTLV()}}
}
