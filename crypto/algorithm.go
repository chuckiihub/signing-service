package crypto

type SignatureAlgorithm int

const (
	SignatureAlgorithmRSA SignatureAlgorithm = iota
	SignatureAlgorithmECC
)

func (s SignatureAlgorithm) String() string {
	switch s {
	case SignatureAlgorithmRSA:
		return "RSA"
	case SignatureAlgorithmECC:
		return "ECC"
	default:
		return "Unkown"
	}
}

func GetSupportedAlgorithms() []string {
	return []string{SignatureAlgorithmECC.String(), SignatureAlgorithmRSA.String()}
}
