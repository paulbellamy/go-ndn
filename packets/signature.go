package packets

type Signature interface {
	Type() uint64
	Bytes() []byte
}

type Sha256WithRSASignature []byte

func (s Sha256WithRSASignature) Type() uint64 {
	return 1
}

func (s Sha256WithRSASignature) Bytes() []byte {
	return []byte(s)
}
