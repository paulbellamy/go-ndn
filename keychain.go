package ndn

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"io"
)

type KeyChain struct{}

func NewKeyChain() *KeyChain {
	return &KeyChain{}
}

type signable interface {
	WriteTo(io.Writer) (int64, error)
	GetName() Name
	SetSignature(Signature)
}

func (k *KeyChain) Sign(packet signable, certificateName Name) error {
	hash := sha256.New()
	_, err := packet.WriteTo(hash)
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

	packet.SetSignature(Sha256WithRSASignature(sig))

	return err
}

// Should ask the KeyLocator which private key to use for this name, but just generating one for now.
func (k *KeyChain) getKey(name Name) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(randSource, 2048)
}
