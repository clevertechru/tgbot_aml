package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

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

type ChainabuseProvider struct {
	client  *http.Client
	apiKey  string
	baseURL string
	logger  *zap.Logger
}

func NewChainabuseProvider() *ChainabuseProvider {
	logger, _ := zap.NewProduction()
	return &ChainabuseProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey:  os.Getenv("AML_API_KEY"),
		baseURL: "https://api.aml-provider.com",
		logger:  logger,
	}
}

type AMLResponse struct {
	Address      string   `json:"address"`
	IsSuspicious bool     `json:"is_suspicious"`
	RiskScore    float64  `json:"risk_score"`
	Details      []string `json:"details"`
}

type TransactionResponse struct {
	TransactionID string   `json:"transaction_id"`
	IsSuspicious  bool     `json:"is_suspicious"`
	RiskScore     float64  `json:"risk_score"`
	Details       []string `json:"details"`
}

func (p *ChainabuseProvider) CheckAddress(ctx context.Context, address string) (*AMLResult, error) {
	if address == "" {
		return nil, ErrEmptyAddress
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/check/address/%s", p.baseURL, address), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp AMLResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &AMLResult{
		Address:      apiResp.Address,
		IsSuspicious: apiResp.IsSuspicious,
		Details:      apiResp.Details,
	}, nil
}

func (p *ChainabuseProvider) CheckTransaction(ctx context.Context, txHash string) (*TransactionResult, error) {
	if txHash == "" {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/check/transaction/%s", p.baseURL, txHash), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &TransactionResult{
		TransactionID: apiResp.TransactionID,
		IsSuspicious:  apiResp.IsSuspicious,
		Details:       apiResp.Details,
	}, nil
}
