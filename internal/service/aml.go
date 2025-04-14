package service

import (
	"github.com/clevertechru/tgbot_aml/internal/domain"
)

// AMLService implements the domain.AMLService interface
type AMLService struct {
	provider domain.AMLService
}

// NewAMLService creates a new instance of AMLService
func NewAMLService(provider domain.AMLService) domain.AMLService {
	return &AMLService{
		provider: provider,
	}
}

// CheckAddress checks if an address is suspicious
func (s *AMLService) CheckAddress(address string) (*domain.CheckResult, error) {
	return s.provider.CheckAddress(address)
}

// CheckTransaction checks if a transaction is suspicious
func (s *AMLService) CheckTransaction(fromAddress, toAddress string, amount float64) (*domain.CheckResult, error) {
	return s.provider.CheckTransaction(fromAddress, toAddress, amount)
}
