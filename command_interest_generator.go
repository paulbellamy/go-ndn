package ndn

import (
	"bytes"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/encoding/tlv"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
)

type commandInterestGenerator struct {
	lastTimestamp time.Time
}

func (c *commandInterestGenerator) generate(i *packets.Interest, keyChain *KeyChain, certificateName name.Name, encoderFactory encoding.EncoderFactory) error {
	timestamp := c.nextTimestamp()

	n := i.GetName()
	n = c.appendTimestamp(n, timestamp)
	n, err := c.appendRandom(n)
	if err != nil {
		return err
	}

	keyChain.Sign(i, certificateName, encoderFactory)

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

func (c *commandInterestGenerator) appendTimestamp(n name.Name, timestamp time.Time) name.Name {
	buf := &bytes.Buffer{}
	tlv.WriteUint(buf, uint64(timestamp.UnixNano()/int64(time.Millisecond)))
	return n.AppendBytes(buf.Bytes())
}

func (c *commandInterestGenerator) appendRandom(n name.Name) (name.Name, error) {
	randomBuffer := make([]byte, 8)
	_, err := randSource.Read(randomBuffer)
	return n.AppendBytes(randomBuffer), err
}

func (c *commandInterestGenerator) setDefaultInterestLifetime(i *packets.Interest) {
	if i.GetInterestLifetime() < 0 {
		// Caller hasn't set this yet, set a default
		i.SetInterestLifetime(1 * time.Second)
	}
}
