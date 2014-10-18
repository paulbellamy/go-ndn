package ndn

import "errors"

var PacketTooLargeError = errors.New("Packet size exceeds maximum limit")
var ErrNameRequired = errors.New("name is required")
var ErrCommandKeyChainNotSet = errors.New("registerPrefix: The command KeyChain has not been set. You must call setCommandSigningInfo.")
var ErrCommandCertificateNameNotSet = errors.New("registerPrefix: The command certificate name has not been set. You must call setCommandSigningInfo.")
var ErrCertificateNotFound = errors.New("certificate not found")
