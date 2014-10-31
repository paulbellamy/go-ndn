package tlv

type Marshaler interface {
	MarshalTLV() (TLV, error)
}
