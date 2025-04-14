package service

import (
	"testing"

	"github.com/clevertechru/tgbot_aml/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct{}

func (m *mockProvider) CheckAddress(address string) (*domain.CheckResult, error) {
	return &domain.CheckResult{
		IsSuspicious: false,
		RiskScore:    0.1,
		Details:      "test",
	}, nil
}

func (m *mockProvider) CheckTransaction(fromAddress, toAddress string, amount float64) (*domain.CheckResult, error) {
	return &domain.CheckResult{
		IsSuspicious: false,
		RiskScore:    0.2,
		Details:      "test",
	}, nil
}

func TestAMLService_CheckAddress(t *testing.T) {
	provider := &mockProvider{}
	service := NewAMLService(provider)

	result, err := service.CheckAddress("test-address")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsSuspicious)
	assert.Equal(t, 0.1, result.RiskScore)
	assert.Equal(t, "test", result.Details)
}

func TestAMLService_CheckTransaction(t *testing.T) {
	provider := &mockProvider{}
	service := NewAMLService(provider)

	result, err := service.CheckTransaction("from", "to", 1.0)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsSuspicious)
	assert.Equal(t, 0.2, result.RiskScore)
	assert.Equal(t, "test", result.Details)
}
