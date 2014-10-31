package ndn

import (
	"io"

	"github.com/paulbellamy/go-ndn/name"
)

type controlParameters struct {
	name name.Name
}

// TODO: Implement this
func (c *controlParameters) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
