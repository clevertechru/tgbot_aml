package domain

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyAddress     = errors.New("address cannot be empty")
	ErrEmptyTransaction = errors.New("transaction hash cannot be empty")
)

type AMLResult struct {
	Address      string
	IsSuspicious bool
	RiskScore    float64
	Details      []string
}

type TransactionResult struct {
	TransactionID string
	IsSuspicious  bool
	RiskScore     float64
	Details       []string
}

// AMLService defines the interface for AML checking services
type AMLService interface {
	// CheckAddress checks if a given address is suspicious
	CheckAddress(address string) (*CheckResult, error)
	// CheckTransaction checks if a transaction is suspicious
	CheckTransaction(fromAddress, toAddress string, amount float64) (*CheckResult, error)
}

// CheckResult represents the result of an AML check
type CheckResult struct {
	IsSuspicious bool
	RiskScore    float64
	Details      string
}

// MockAMLService is a mock implementation of AMLService for testing
type MockAMLService struct{}

// NewMockAMLService creates a new instance of MockAMLService
func NewMockAMLService() *MockAMLService {
	return &MockAMLService{}
}

// CheckAddress implements AMLService interface
func (m *MockAMLService) CheckAddress(address string) (*CheckResult, error) {
	return &CheckResult{
		IsSuspicious: false,
		RiskScore:    0.1,
		Details:      "Address appears to be clean",
	}, nil
}

// CheckTransaction implements AMLService interface
func (m *MockAMLService) CheckTransaction(fromAddress, toAddress string, amount float64) (*CheckResult, error) {
	return &CheckResult{
		IsSuspicious: false,
		RiskScore:    0.2,
		Details:      "Transaction appears to be clean",
	}, nil
}

func FormatAMLResult(result *AMLResult) string {
	return fmt.Sprintf("Address %s checked against AML databases:\nRisk Score: %.2f\n- %s",
		result.Address,
		result.RiskScore,
		result.Details[0])
}

func FormatTransactionResult(result *TransactionResult) string {
	return fmt.Sprintf("Transaction %s checked against AML databases:\nRisk Score: %.2f\n- %s",
		result.TransactionID,
		result.RiskScore,
		result.Details[0])
}
