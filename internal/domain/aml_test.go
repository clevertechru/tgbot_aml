package domain

import (
	"testing"
)

type MockAMLService struct{}

func NewMockAMLService() *MockAMLService {
	return &MockAMLService{}
}

func (m *MockAMLService) CheckAddress(address string) (*CheckResult, error) {
	return &CheckResult{
		IsSuspicious: false,
		RiskScore:    0.1,
		Details:      "Address appears to be clean",
	}, nil
}

func (m *MockAMLService) CheckTransaction(fromAddress, toAddress string, amount float64) (*CheckResult, error) {
	return &CheckResult{
		IsSuspicious: false,
		RiskScore:    0.2,
		Details:      "Transaction appears to be clean",
	}, nil
}

func TestMockAMLService_CheckAddress(t *testing.T) {
	service := NewMockAMLService()
	result, err := service.CheckAddress("test-address")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsSuspicious {
		t.Error("expected address to be clean")
	}
	if result.RiskScore != 0.1 {
		t.Errorf("expected risk score 0.1, got %f", result.RiskScore)
	}
}

func TestMockAMLService_CheckTransaction(t *testing.T) {
	service := NewMockAMLService()
	result, err := service.CheckTransaction("from", "to", 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsSuspicious {
		t.Error("expected transaction to be clean")
	}
	if result.RiskScore != 0.2 {
		t.Errorf("expected risk score 0.2, got %f", result.RiskScore)
	}
}
