package domain

import (
	"github.com/chuckiihub/signing-service/crypto"
)

type Device struct {
	UUID             string                    `json:"uuid"`
	Label            string                    `json:"label"`
	SignatureCounter int                       `json:"signatureCounter"`
	Algorithm        crypto.SignatureAlgorithm `json:"algorithm"`
	PublicKey        []byte                    `json:"publicKey"`
	PrivateKey       []byte                    `json:"privateKey"`
	LastSignature    string                    `json:"lastSignature"`
}
