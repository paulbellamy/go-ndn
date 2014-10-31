package packets

import "github.com/paulbellamy/go-ndn/name"

type Data struct {
	name name.Name
	MetaInfo
	content   []byte
	signature Signature
}

func (d *Data) GetName() name.Name {
	if d.name == nil {
		d.name = name.Name{}
	}
	return d.name
}

func (d *Data) SetName(x name.Name) {
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
