package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTranslations(t *testing.T) {
	translations := GetTranslations()
	assert.NotNil(t, translations)
	assert.NotEmpty(t, translations.WelcomeMessage)
	assert.NotEmpty(t, translations.CheckUsage)
	assert.NotEmpty(t, translations.CheckTxUsage)
	assert.NotEmpty(t, translations.UnknownCommand)
	assert.NotEmpty(t, translations.AddressChecked)
	assert.NotEmpty(t, translations.ErrorChecking)
	assert.NotEmpty(t, translations.LanguageSelection)
}

func TestFormatAddressChecked(t *testing.T) {
	expected := "Address 0x1234 checked against AML databases:\n- No suspicious activity found"
	actual := FormatAddressChecked("0x1234")
	assert.Equal(t, expected, actual)
}

func TestFormatErrorChecking(t *testing.T) {
	expected := "Error checking address: test error"
	actual := FormatErrorChecking("test error")
	assert.Equal(t, expected, actual)
}

func TestLoadTranslations(t *testing.T) {
	// Test loading English translations
	translations := loadTranslations("en")
	assert.NotNil(t, translations)
	assert.NotEmpty(t, translations.WelcomeMessage)
	assert.NotEmpty(t, translations.CheckUsage)
	assert.NotEmpty(t, translations.CheckTxUsage)
	assert.NotEmpty(t, translations.UnknownCommand)
	assert.NotEmpty(t, translations.AddressChecked)
	assert.NotEmpty(t, translations.ErrorChecking)
	assert.NotEmpty(t, translations.LanguageSelection)

	// Test loading non-existent language (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when loading non-existent language")
		}
	}()
	loadTranslations("fr")
}
