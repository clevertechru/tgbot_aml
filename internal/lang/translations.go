package lang

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Language string

const (
	English Language = "en"
	Russian Language = "ru"
)

var (
	translations    = make(map[Language]map[string]string)
	translationsDir string
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	translationsDir = filepath.Join(filepath.Dir(filename), "translations")
	loadTranslations()
}

func loadTranslations() {
	// Load English translations
	enPath := filepath.Join(translationsDir, "en.yml")
	enTranslations, err := loadYAMLFile(enPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load English translations: %v", err))
	}
	translations[English] = enTranslations

	// Load Russian translations
	ruPath := filepath.Join(translationsDir, "ru.yml")
	ruTranslations, err := loadYAMLFile(ruPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load Russian translations: %v", err))
	}
	translations[Russian] = ruTranslations
}

func loadYAMLFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var translations map[string]string
	if err := yaml.Unmarshal(data, &translations); err != nil {
		return nil, err
	}

	return translations, nil
}

func Get(lang Language, key string, args ...interface{}) string {
	if lang == "" {
		lang = English // default to English
	}

	msg, ok := translations[lang][key]
	if !ok {
		msg = translations[English][key] // fallback to English
	}

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}
