package crypto

import (
	"testing"
)

func supportedAlgorithms() []SignatureAlgorithm {
	return []SignatureAlgorithm{SignatureAlgorithmRSA, SignatureAlgorithmECC}
}

func createCryptoAndKeyPair(t *testing.T, algorithm SignatureAlgorithm) (Crypto, KeyPair) {
	crypto, err := NewCrypto(algorithm)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto: %v", err)
	}

	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	return crypto, keyPair
}

func TestRSAKeyPairGeneration(t *testing.T) {
	_, keyPair := createCryptoAndKeyPair(t, SignatureAlgorithmRSA)

	rsaKeyPair, ok := keyPair.(*RSAKeyPair)
	if !ok {
		t.Fatalf("Generated key pair is not of type %vKeyPair", SignatureAlgorithmRSA)
	}

	if rsaKeyPair.Public == nil {
		t.Errorf("%v public key is nil", SignatureAlgorithmRSA)
	}

	if rsaKeyPair.Private == nil {
		t.Errorf("%v private key is nil", SignatureAlgorithmRSA)
	}
}

func TestECCKeyPairGeneration(t *testing.T) {
	_, keyPair := createCryptoAndKeyPair(t, SignatureAlgorithmECC)

	eccKeyPair, ok := keyPair.(*ECCKeyPair)
	if !ok {
		t.Fatalf("Generated key pair is not of type %vKeyPair", SignatureAlgorithmECC)
	}

	if eccKeyPair.Public == nil {
		t.Errorf("%v public key is nil", SignatureAlgorithmECC)
	}

	if eccKeyPair.Private == nil {
		t.Errorf("%v private key is nil", SignatureAlgorithmECC)
	}
}

func TestRSASignatureCreationAndValidation(t *testing.T) {
	const TEST_DATA = "The Beatles Are Great, specially the White Album"

	for _, algorithm := range supportedAlgorithms() {
		crypto, keyPair := createCryptoAndKeyPair(t, algorithm)

		_, privateKey, err := crypto.Marshal(keyPair)
		if err != nil {
			t.Fatalf("Failed to marshal %v key pair: %v", algorithm, err)
		}

		data := []byte(TEST_DATA)
		signature, err := crypto.Sign(data, privateKey)
		if err != nil {
			t.Fatalf("Failed to create RSA signature: %v", err)
		}

		valid, err := crypto.Verify(data, signature, privateKey)

		if !valid || err != nil {
			t.Errorf("%v signature verification failed %v", algorithm, err)
		}
	}
}

func TestRSASignatureCreationAndValidationFailureWithTamperedData(t *testing.T) {
	const TEST_DATA = "The Beatles Are Great, specially the White Album"
	const TAMPERED_DATA = "The Beatles Are Not Great, specially Rubber Soul."

	for _, algorithm := range supportedAlgorithms() {
		crypto, err := NewCrypto(algorithm)
		if err != nil {
			t.Fatalf("Failed to create %v crypto: %v", algorithm, err)
		}

		keyPair, err := crypto.GenerateKeyPair()
		if err != nil {
			t.Fatalf("Failed to generate %v key pair: %v", algorithm, err)
		}

		publicKey, privateKey, err := crypto.Marshal(keyPair)
		if err != nil {
			t.Fatalf("Failed to marshal %v key pair: %v", algorithm, err)
		}

		signature, err := crypto.Sign([]byte(TEST_DATA), privateKey)
		if err != nil {
			t.Fatalf("Failed to create %v signature: %v", algorithm, err)
		}

		valid, err := crypto.Verify([]byte(TAMPERED_DATA), signature, publicKey)

		if valid || err == nil {
			t.Errorf("%v signature verification should have failed with tampered data", algorithm)
		}
	}
}

func TestRSAMarshalUnmarshal(t *testing.T) {
	crypto, keyPair := createCryptoAndKeyPair(t, SignatureAlgorithmRSA)
	rsaKeyPair, ok := keyPair.(*RSAKeyPair)
	if !ok {
		t.Fatal("Generated key pair is not of type RSAKeyPair")
	}

	_, privateKey, err := crypto.Marshal(keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal %s key pair: %v", SignatureAlgorithmRSA, err)
	}

	unmarshalledKeyPair, err := crypto.Unmarshal(privateKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal %s key pair: %v", SignatureAlgorithmRSA, err)
	}
	unmarshalledRsaKeyPair, ok := unmarshalledKeyPair.(*RSAKeyPair)
	if !ok {
		t.Fatal("Unmarshalled key pair is not of type RSAKeyPair")
	}

	if !unmarshalledRsaKeyPair.Public.Equal(rsaKeyPair.Public) {
		t.Errorf("%s public key does not match after marshal and unmarshal", SignatureAlgorithmRSA)
	}

	if !unmarshalledRsaKeyPair.Private.Equal(rsaKeyPair.Private) {
		t.Errorf("%s public key does not match after marshal and unmarshal", SignatureAlgorithmRSA)
	}
}

func TestECCMarshalUnmarshal(t *testing.T) {
	crypto, keyPair := createCryptoAndKeyPair(t, SignatureAlgorithmECC)
	eccKeyPair, ok := keyPair.(*ECCKeyPair)
	if !ok {
		t.Fatal("Generated key pair is not of type ECCKeyPair")
	}

	_, privateKey, err := crypto.Marshal(keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal %s key pair: %v", SignatureAlgorithmECC, err)
	}

	unmarshalledKeyPair, err := crypto.Unmarshal(privateKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal %s key pair: %v", SignatureAlgorithmECC, err)
	}
	unmarshalledEccKeyPair, ok := unmarshalledKeyPair.(*ECCKeyPair)
	if !ok {
		t.Fatal("Unmarshalled key pair is not of type ECCKeyPair")
	}

	if !unmarshalledEccKeyPair.Public.Equal(eccKeyPair.Public) {
		t.Errorf("%s public key does not match after marshal and unmarshal", SignatureAlgorithmECC)
	}

	if !unmarshalledEccKeyPair.Private.Equal(eccKeyPair.Private) {
		t.Errorf("%s public key does not match after marshal and unmarshal", SignatureAlgorithmECC)
	}
}
