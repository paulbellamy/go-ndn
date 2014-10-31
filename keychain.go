package ndn

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
)

type KeyChain struct{}

func NewKeyChain() *KeyChain {
	return &KeyChain{}
}

type signable interface {
	GetName() name.Name
	SetSignature(packets.Signature)
}

func (k *KeyChain) Sign(packet signable, certificateName name.Name, newEncoder encoding.EncoderFactory) error {
	hash := sha256.New()
	err := newEncoder(hash).Encode(packet)
	if err != nil {
		return err
	}

	privateKey, err := k.getKey(packet.GetName())
	if err != nil {
		return err
	}

	sig, err := rsa.SignPKCS1v15(randSource, privateKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return err
	}

	packet.SetSignature(packets.Sha256WithRSASignature(sig))

	return err
}

// Should ask the KeyLocator which private key to use for this name, but just generating one for now.
func (k *KeyChain) getKey(n name.Name) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(randSource, 2048)
}
