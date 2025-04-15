package services

import (
	"context"
	"errors"
	"testing"

	"github.com/clevertechru/tgbot_aml/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockAMLProvider struct {
	mock.Mock
}

func (m *MockAMLProvider) CheckAddress(ctx context.Context, address string) (*domain.CheckResult, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CheckResult), args.Error(1)
}

func (m *MockAMLProvider) CheckTransaction(ctx context.Context, txHash string) (*domain.CheckResult, error) {
	args := m.Called(ctx, txHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CheckResult), args.Error(1)
}

type AMLServiceTestSuite struct {
	suite.Suite
	provider *MockAMLProvider
	service  *AMLService
	ctx      context.Context
}

func (s *AMLServiceTestSuite) SetupTest() {
	s.provider = new(MockAMLProvider)
	s.service = NewAMLService(s.provider)
	s.ctx = context.Background()
}

func TestAMLServiceSuite(t *testing.T) {
	suite.Run(t, new(AMLServiceTestSuite))
}

func (s *AMLServiceTestSuite) TestCheckAddress() {
	testCases := []struct {
		name        string
		address     string
		mockResult  *domain.CheckResult
		mockError   error
		expectError bool
		expected    *domain.AMLResult
	}{
		{
			name:    "Clean address",
			address: "0x123",
			mockResult: &domain.CheckResult{
				IsSuspicious: false,
				RiskScore:    0.1,
				Details:      "No suspicious activity",
			},
			expected: &domain.AMLResult{
				Address:      "0x123",
				IsSuspicious: false,
				RiskScore:    0.1,
				Details:      []string{"No suspicious activity"},
			},
		},
		{
			name:    "Suspicious address",
			address: "0x456",
			mockResult: &domain.CheckResult{
				IsSuspicious: true,
				RiskScore:    0.8,
				Details:      "High-risk activity detected",
			},
			expected: &domain.AMLResult{
				Address:      "0x456",
				IsSuspicious: true,
				RiskScore:    0.8,
				Details:      []string{"High-risk activity detected"},
			},
		},
		{
			name:        "Provider error",
			address:     "0x789",
			mockError:   errors.New("provider error"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.provider.On("CheckAddress", s.ctx, tc.address).Return(tc.mockResult, tc.mockError).Once()

			result, err := s.service.CheckAddress(s.ctx, tc.address)

			if tc.expectError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tc.expected, result)
			}

			s.provider.AssertExpectations(s.T())
		})
	}
}

func (s *AMLServiceTestSuite) TestCheckTransaction() {
	testCases := []struct {
		name        string
		txHash      string
		mockResult  *domain.CheckResult
		mockError   error
		expectError bool
		expected    *domain.TransactionResult
	}{
		{
			name:   "Clean transaction",
			txHash: "0xabc",
			mockResult: &domain.CheckResult{
				IsSuspicious: false,
				RiskScore:    0.2,
				Details:      "Normal transaction",
			},
			expected: &domain.TransactionResult{
				TransactionID: "0xabc",
				IsSuspicious:  false,
				RiskScore:     0.2,
				Details:       []string{"Normal transaction"},
			},
		},
		{
			name:   "Suspicious transaction",
			txHash: "0xdef",
			mockResult: &domain.CheckResult{
				IsSuspicious: true,
				RiskScore:    0.9,
				Details:      "Suspicious pattern detected",
			},
			expected: &domain.TransactionResult{
				TransactionID: "0xdef",
				IsSuspicious:  true,
				RiskScore:     0.9,
				Details:       []string{"Suspicious pattern detected"},
			},
		},
		{
			name:        "Provider error",
			txHash:      "0xghi",
			mockError:   errors.New("provider error"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.provider.On("CheckTransaction", s.ctx, tc.txHash).Return(tc.mockResult, tc.mockError).Once()

			result, err := s.service.CheckTransaction(s.ctx, tc.txHash)

			if tc.expectError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tc.expected, result)
			}

			s.provider.AssertExpectations(s.T())
		})
	}
}
