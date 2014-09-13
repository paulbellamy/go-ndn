package encoding

import "io"

func WriteTLV(w io.Writer, tlv *TLV) error {
	err := WriteNumber(w, tlv.Type)
	if err != nil {
		return err
	}

	err = WriteNumber(w, uint64(len(tlv.Value)))
	if err != nil {
		return err
	}

	_, err = w.Write(tlv.Value)
	return err
}

func (t *TLV) WriteTo(w io.Writer) error {
	return WriteTLV(w, t)
}
