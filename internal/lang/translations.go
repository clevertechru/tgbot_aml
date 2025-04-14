package lang

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Translations holds all translation strings
type Translations struct {
	WelcomeMessage    string `yaml:"welcome_message"`
	CheckUsage        string `yaml:"check_usage"`
	CheckTxUsage      string `yaml:"check_tx_usage"`
	UnknownCommand    string `yaml:"unknown_command"`
	AddressChecked    string `yaml:"address_checked"`
	ErrorChecking     string `yaml:"error_checking"`
	LanguageSelection string `yaml:"language_selection"`
}

var translations *Translations

func init() {
	translations = loadTranslations("en")
}

// GetTranslations returns the loaded translations
func GetTranslations() *Translations {
	return translations
}

func loadTranslations(lang string) *Translations {
	filename := filepath.Join("internal", "lang", "translations", fmt.Sprintf("%s.yml", lang))
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to load %s translations: %v", lang, err))
	}

	var t Translations
	if err := yaml.Unmarshal(data, &t); err != nil {
		panic(fmt.Sprintf("Failed to parse %s translations: %v", lang, err))
	}

	return &t
}

// FormatAddressChecked formats the address checked message
func FormatAddressChecked(address string) string {
	return fmt.Sprintf(translations.AddressChecked, address)
}

// FormatErrorChecking formats the error checking message
func FormatErrorChecking(err string) string {
	return fmt.Sprintf(translations.ErrorChecking, err)
}
