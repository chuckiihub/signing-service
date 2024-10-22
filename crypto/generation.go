package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

// To support new algorithms easily, I've implemented a KeyPair structure that can be
// be used for algorithms that have a PublicKey and a PrivateKey. Then, you can use the Crypto interface
// to generate, marshal and unmarshal, signing and verifying signatures using the mentioned KeyPair

// NewKeyGenerator creates a new KeyGenerator for the given signature algorithm.
func NewKeyGenerator(s SignatureAlgorithm) (KeyGenerator, error) {
	switch s {
	case SignatureAlgorithmRSA:
		return &RSAGenerator{}, nil
	case SignatureAlgorithmECC:
		return &ECCGenerator{}, nil
	default:
		return nil, errors.New("not implemented algorithm")
	}
}

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (KeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

func (g *RSAGenerator) Marshal(keyPair KeyPair) ([]byte, []byte, error) {
	rsaKeyPair, ok := keyPair.(*RSAKeyPair)
	if !ok {
		return nil, nil, errors.New("invalid key pair")
	}
	marshaller := NewRSAMarshaler()
	return marshaller.Marshal(*rsaKeyPair)
}

func (g *RSAGenerator) Unmarshal(privateKey []byte) (KeyPair, error) {
	marshaller := NewRSAMarshaler()
	return marshaller.Unmarshal(privateKey)
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (KeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

func (g *ECCGenerator) Marshal(keyPair KeyPair) ([]byte, []byte, error) {
	eccKeyPair, ok := keyPair.(*ECCKeyPair)
	if !ok {
		return nil, nil, errors.New("invalid key pair")
	}

	marshaller := NewECCMarshaler()
	return marshaller.Encode(*eccKeyPair)
}

func (g *ECCGenerator) Unmarshal(privateKey []byte) (KeyPair, error) {
	marshaller := NewECCMarshaler()
	return marshaller.Decode(privateKey)
}
