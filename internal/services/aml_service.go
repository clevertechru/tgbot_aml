package services

import (
	"context"

	"github.com/clevertechru/tgbot_aml/internal/domain"
)

type AMLService struct {
	provider interface {
		CheckAddress(ctx context.Context, address string) (*domain.CheckResult, error)
		CheckTransaction(ctx context.Context, txHash string) (*domain.CheckResult, error)
	}
}

func NewAMLService(provider interface {
	CheckAddress(ctx context.Context, address string) (*domain.CheckResult, error)
	CheckTransaction(ctx context.Context, txHash string) (*domain.CheckResult, error)
}) *AMLService {
	return &AMLService{
		provider: provider,
	}
}

func (s *AMLService) CheckAddress(ctx context.Context, address string) (*domain.AMLResult, error) {
	result, err := s.provider.CheckAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	return &domain.AMLResult{
		Address:      address,
		IsSuspicious: result.IsSuspicious,
		RiskScore:    result.RiskScore,
		Details:      []string{result.Details},
	}, nil
}

func (s *AMLService) CheckTransaction(ctx context.Context, txHash string) (*domain.TransactionResult, error) {
	result, err := s.provider.CheckTransaction(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return &domain.TransactionResult{
		TransactionID: txHash,
		IsSuspicious:  result.IsSuspicious,
		RiskScore:     result.RiskScore,
		Details:       []string{result.Details},
	}, nil
}
