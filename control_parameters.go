package ndn

import "io"

type controlParameters struct {
	name Name
}

// TODO: Implement this
func (c *controlParameters) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
