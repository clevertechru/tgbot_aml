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

// FormatAMLResult formats an AMLResult into a human-readable string
func FormatAMLResult(result *AMLResult) string {
	if result == nil {
		return "No result available"
	}

	status := "✅ Clean"
	if result.IsSuspicious {
		status = "⚠️ Suspicious"
	}

	return fmt.Sprintf(
		"Address: %s\nStatus: %s\nRisk Score: %.2f\nDetails:\n%s",
		result.Address,
		status,
		result.RiskScore,
		formatDetails(result.Details),
	)
}

// FormatTransactionResult formats a TransactionResult into a human-readable string
func FormatTransactionResult(result *TransactionResult) string {
	if result == nil {
		return "No result available"
	}

	status := "✅ Clean"
	if result.IsSuspicious {
		status = "⚠️ Suspicious"
	}

	return fmt.Sprintf(
		"Transaction: %s\nStatus: %s\nRisk Score: %.2f\nDetails:\n%s",
		result.TransactionID,
		status,
		result.RiskScore,
		formatDetails(result.Details),
	)
}

func formatDetails(details []string) string {
	if len(details) == 0 {
		return "  - No details available"
	}

	var formatted string
	for _, detail := range details {
		formatted += fmt.Sprintf("  - %s\n", detail)
	}
	return formatted
}
