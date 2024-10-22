package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

// Tries to make an attemp to support a variable types of encryption algorithms.
// I would refactor this interface in the future to avoid using empty interfaces.
type KeyPair interface {
	PublicKey() interface{}
	PrivateKey() interface{}
}

type KeyGenerator interface {
	Generate() (KeyPair, error)
	Marshal(keyPair KeyPair) ([]byte, []byte, error)
	Unmarshal(privateKey []byte) (KeyPair, error)
}

// If a new algorithm is needed, the algorithm should implement this interface.
// It should also register itself in the NewCrypto function and add the algorithm name
// to the SignatureAlgorithm enum.
// This interface combines the functionality 2 interfaces: KeyGenerator and Signer.
// The challenge wanted me to code a Signer interface so I've decided to add a Generator and implement
// the Signer interfaces separated (and not to create just one Crypto interface that implements both).
type Crypto interface {
	GenerateKeyPair() (KeyPair, error)
	Sign(dataToBeSigned []byte, privateKey []byte) ([]byte, error)
	Verify(dataToBeSigned []byte, signature []byte, privateKey []byte) (bool, error)
	Marshal(keyPair KeyPair) ([]byte, []byte, error)
	Unmarshal(privateKey []byte) (KeyPair, error)
}

func NewCrypto(algorithm SignatureAlgorithm) (Crypto, error) {
	switch algorithm {
	case SignatureAlgorithmRSA:
		return &RSACrypto{}, nil
	case SignatureAlgorithmECC:
		return &ECCCrypto{}, nil
	}

	return nil, errors.New("not implemented algorithm")
}

type RSACrypto struct{}

func (c *RSACrypto) GenerateKeyPair() (KeyPair, error) {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

func (c *RSACrypto) Verify(dataToBeSigned []byte, signature []byte, publicKey []byte) (bool, error) {
	signer, err := CreateSigner(SignatureAlgorithmRSA, publicKey)
	if err != nil {
		return false, err
	}

	return signer.Verify(dataToBeSigned, signature), nil
}

func (c *RSACrypto) Sign(dataToBeSigned []byte, privateKey []byte) ([]byte, error) {
	signer, err := CreateSigner(SignatureAlgorithmRSA, privateKey)
	if err != nil {
		return nil, err
	}

	return signer.Sign(dataToBeSigned)
}

func (c *RSACrypto) Marshal(keyPair KeyPair) ([]byte, []byte, error) {
	generator, err := NewKeyGenerator(SignatureAlgorithmRSA)
	if err != nil {
		return nil, nil, err
	}

	return generator.Marshal(keyPair)
}

func (c *RSACrypto) Unmarshal(privateKey []byte) (KeyPair, error) {
	generator, err := NewKeyGenerator(SignatureAlgorithmRSA)
	if err != nil {
		return nil, err
	}

	return generator.Unmarshal(privateKey)
}

type ECCCrypto struct{}

func (c *ECCCrypto) GenerateKeyPair() (KeyPair, error) {
	generator, err := NewKeyGenerator(SignatureAlgorithmECC)
	if err != nil {
		return nil, err
	}

	return generator.Generate()
}

func (c *ECCCrypto) Verify(dataToBeSigned []byte, signature []byte, privateKey []byte) (bool, error) {
	signer, err := CreateSigner(SignatureAlgorithmECC, privateKey)
	if err != nil {
		return false, err
	}

	return signer.Verify(dataToBeSigned, signature), nil
}

func (c *ECCCrypto) Sign(dataToBeSigned []byte, privateKey []byte) ([]byte, error) {
	signer, err := CreateSigner(SignatureAlgorithmECC, privateKey)
	if err != nil {
		return nil, err
	}

	return signer.Sign(dataToBeSigned)
}

func (c *ECCCrypto) Marshal(keyPair KeyPair) ([]byte, []byte, error) {
	generator, err := NewKeyGenerator(SignatureAlgorithmECC)
	if err != nil {
		return nil, nil, err
	}

	return generator.Marshal(keyPair)
}

func (c *ECCCrypto) Unmarshal(privateKey []byte) (KeyPair, error) {
	generator, err := NewKeyGenerator(SignatureAlgorithmECC)
	if err != nil {
		return nil, err
	}

	return generator.Unmarshal(privateKey)
}
