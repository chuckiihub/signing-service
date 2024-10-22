package domain

type PersistenceHealth struct {
	// in the future, we can add more fields to describe the health of the persistence layer
	// we could even add latency metrics here for them
	Status string `json:"status"`
}

type ServiceHealth struct {
	Status           string                       `json:"status"`
	PersistenceLayer map[string]PersistenceHealth `json:"persitence_layer"`
}

type Health struct {
	Status   string                   `json:"status"`
	Version  string                   `json:"version"`
	Services map[string]ServiceHealth `json:"services"`
}

const (
	HealthStatusPass   = "pass"
	HealthStatusFailed = "failed"
)
