package packets

type Signature interface {
	Type() uint64
	Bytes() []byte
}

type DigestSha256 []byte

func (s DigestSha256) Type() uint64 {
	return 0
}

func (s DigestSha256) Bytes() []byte {
	return []byte(s)
}

type Sha256WithRSASignature []byte

func (s Sha256WithRSASignature) Type() uint64 {
	return 1
}

func (s Sha256WithRSASignature) Bytes() []byte {
	return []byte(s)
}

type Sha256WithECDSASignature []byte

func (s Sha256WithECDSASignature) Type() uint64 {
	return 3
}

func (s Sha256WithECDSASignature) Bytes() []byte {
	return []byte(s)
}
