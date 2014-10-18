package ndn

import "errors"

var PacketTooLargeError = errors.New("Packet size exceeds maximum limit")
var ErrNameRequired = errors.New("name is required")
