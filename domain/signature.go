package domain

type Signature struct {
	UUID       string `json:"uuid"`
	DeviceUUID string `json:"deviceId"`
	SignedData string `json:"signedData"`
	Signature  string `json:"signature"`
}
