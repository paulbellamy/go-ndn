package ndn

import (
	"bytes"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
)

type commandInterestGenerator struct {
	lastTimestamp time.Time
}

func (c *commandInterestGenerator) generate(i *Interest, keyChain *KeyChain, certificateName Name) error {
	timestamp := c.nextTimestamp()

	name := i.GetName()
	name = c.appendTimestamp(name, timestamp)
	name, err := c.appendRandom(name)
	if err != nil {
		return err
	}

	keyChain.Sign(i, certificateName)

	c.setDefaultInterestLifetime(i)

	// We successfully signed the interest, so update the timestamp.
	c.lastTimestamp = timestamp

	return nil
}

// Ensure that the timestamp is unique and monotonic
func (c *commandInterestGenerator) nextTimestamp() time.Time {
	t := time.Now().Truncate(time.Millisecond)
	if !c.lastTimestamp.IsZero() && (t.Before(c.lastTimestamp) || t.Equal(c.lastTimestamp)) {
		t = c.lastTimestamp.Add(1 * time.Millisecond)
	}
	return t
}

func (c *commandInterestGenerator) appendTimestamp(name Name, timestamp time.Time) Name {
	buf := &bytes.Buffer{}
	encoding.WriteUint(buf, uint64(timestamp.UnixNano()/int64(time.Millisecond)))
	return name.AppendBytes(buf.Bytes())
}

func (c *commandInterestGenerator) appendRandom(name Name) (Name, error) {
	randomBuffer := make([]byte, 8)
	_, err := randSource.Read(randomBuffer)
	return name.AppendBytes(randomBuffer), err
}

func (c *commandInterestGenerator) setDefaultInterestLifetime(i *Interest) {
	if i.GetInterestLifetime() < 0 {
		// Caller hasn't set this yet, set a default
		i.SetInterestLifetime(1 * time.Second)
	}
}
