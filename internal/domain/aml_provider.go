package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (p *ChainabuseProvider) SetAPIKey(apiKey string) {
	p.apiKey = apiKey
}

func (p *ChainabuseProvider) SetBaseURL(baseURL string) {
	p.baseURL = baseURL
}

func (p *ChainabuseProvider) CheckAddress(ctx context.Context, address string) (*CheckResult, error) {
	if address == "" {
		return nil, ErrEmptyAddress
	}

	resp, err := p.client.Get(fmt.Sprintf("%s/address/%s", p.baseURL, address))
	if err != nil {
		return nil, fmt.Errorf("failed to check address: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		IsSuspicious bool    `json:"is_suspicious"`
		RiskScore    float64 `json:"risk_score"`
		Details      string  `json:"details"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &CheckResult{
		IsSuspicious: result.IsSuspicious,
		RiskScore:    result.RiskScore,
		Details:      result.Details,
	}, nil
}

func (p *ChainabuseProvider) CheckTransaction(ctx context.Context, txHash string) (*CheckResult, error) {
	if txHash == "" {
		return nil, ErrEmptyTransaction
	}

	resp, err := p.client.Get(fmt.Sprintf("%s/transaction/%s", p.baseURL, txHash))
	if err != nil {
		return nil, fmt.Errorf("failed to check transaction: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		IsSuspicious bool    `json:"is_suspicious"`
		RiskScore    float64 `json:"risk_score"`
		Details      string  `json:"details"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &CheckResult{
		IsSuspicious: result.IsSuspicious,
		RiskScore:    result.RiskScore,
		Details:      result.Details,
	}, nil
}
