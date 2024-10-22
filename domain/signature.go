package domain

type Signature struct {
	DeviceUUID string `json:"deviceId"`
	SignedData string `json:"signedData"`
	Signature  string `json:"signature"`
}
