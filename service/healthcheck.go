package service

import "github.com/chuckiihub/signing-service/domain"

type ServiceHealthCheck interface {
	CheckHealth() domain.ServiceHealth
}
