// Package service implements the core business logic.
package service

import (
	"github.com/elazarl/goproxy"
)

// CAService defines the interface for managing the MITM Certificate Authority.
type CAService interface {
	GetCACert() []byte
}

type caService struct{}

// NewCAService creates a new CAService.
func NewCAService() CAService {
	return &caService{}
}

func (s *caService) GetCACert() []byte {
	return goproxy.CA_CERT
}
