package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"strconv"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	Verify(dataToBeSigned []byte, signature []byte) bool
}

func CreateSigner(t SignatureAlgorithm, privateKey []byte) (Signer, error) {
	switch t {
	case SignatureAlgorithmECC:
		keyPair, err := NewECCMarshaler().Decode(privateKey)
		if err != nil {
			return nil, err
		}

		return NewECCSigner(*keyPair), nil
	case SignatureAlgorithmRSA:
		marshaller := NewRSAMarshaler()

		keyPair, err := marshaller.Unmarshal(privateKey)
		if err != nil {
			return nil, err
		}

		return NewRSASigner(*keyPair), nil
	default:
		return nil, errors.New(`signature algorithm ` + strconv.Itoa(int(t)) + `not implemented.`)
	}
}

type RSASigner struct {
	keyPair RSAKeyPair
}

type ECCSigner struct {
	keyPair ECCKeyPair
}

func NewECCSigner(keyPair ECCKeyPair) ECCSigner {
	return ECCSigner{
		keyPair: keyPair,
	}
}

func NewRSASigner(keyPair RSAKeyPair) RSASigner {
	return RSASigner{
		keyPair: keyPair,
	}
}

func (signer RSASigner) Verify(dataToBeSigned []byte, signature []byte) bool {
	hashed := sha256.Sum256(dataToBeSigned)

	return rsa.VerifyPKCS1v15(signer.keyPair.Public, crypto.SHA256, hashed[:], signature) == nil
}

func (signer RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)

	return rsa.SignPKCS1v15(rand.Reader, signer.keyPair.Private, crypto.SHA256, hashed[:])
}

func (signer ECCSigner) Verify(dataToBeSigned []byte, signature []byte) bool {
	hashed := sha256.Sum256(dataToBeSigned)

	return ecdsa.VerifyASN1(signer.keyPair.Public, hashed[:], signature)
}

func (signer ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	bytes, err := ecdsa.SignASN1(rand.Reader, signer.keyPair.Private, hashed[:])

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
