package encoding

import (
	"fmt"
	"reflect"
)

type UnsupportedTypeError struct {
	Encoding string
	Type     reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("%s: unsupported type: %s", e.Encoding, e.Type.String())
}

type NameRequiredError struct{}

func (e *NameRequiredError) Error() string {
	return "name required"
}

type SignatureRequiredError struct{}

func (e *SignatureRequiredError) Error() string {
	return "signature required"
}

type PacketTooLargeError struct{}

func (e *PacketTooLargeError) Error() string {
	return fmt.Sprintf("packet size exceeds maximum limit of %d bytes", MaxNDNPacketSize)
}

type InvalidUnmarshalError struct {
	Message string
}

func (e *InvalidUnmarshalError) Error() string {
	return e.Message
}

type InvalidUnmarshalTargetError struct {
	Encoding string
	Type     reflect.Type
}

func (e *InvalidUnmarshalTargetError) Error() string {
	if e.Type == nil {
		return e.Encoding + ": Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return e.Encoding + ": Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return e.Encoding + ": Unmarshal(nil " + e.Type.String() + ")"
}
