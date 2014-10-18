package ndn

import "crypto/rand"

const (
	MaxNDNPacketSize = 8800
)

var randSource = rand.Reader
