package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockAMLService_CheckAddress(t *testing.T) {
	service := NewMockAMLService()

	result, err := service.CheckAddress("test-address")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsSuspicious)
	assert.Equal(t, 0.1, result.RiskScore)
	assert.Equal(t, "Address appears to be clean", result.Details)
}

func TestMockAMLService_CheckTransaction(t *testing.T) {
	service := NewMockAMLService()

	result, err := service.CheckTransaction("from", "to", 1.0)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsSuspicious)
	assert.Equal(t, 0.2, result.RiskScore)
	assert.Equal(t, "Transaction appears to be clean", result.Details)
}
