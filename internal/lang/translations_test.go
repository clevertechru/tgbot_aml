package lang

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TranslationsTestSuite struct {
	suite.Suite
	tempDir     string
	originalDir string
}

func (s *TranslationsTestSuite) SetupTest() {
	s.tempDir = s.T().TempDir()
	s.originalDir = translationsDir
	translationsDir = s.tempDir
}

func (s *TranslationsTestSuite) TearDownTest() {
	translationsDir = s.originalDir
	translations = make(map[Language]map[string]string)
}

func TestTranslationsSuite(t *testing.T) {
	suite.Run(t, new(TranslationsTestSuite))
}

func (s *TranslationsTestSuite) TestLoadTranslations() {
	// Create test YAML files
	enContent := `welcome: "Test welcome"
check_usage: "Test check usage"
unknown_command: "Test unknown command"`

	ruContent := `welcome: "Тестовое приветствие"
check_usage: "Тестовая инструкция"
unknown_command: "Тестовая неизвестная команда"`

	// Write test files
	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "en.yml"), []byte(enContent), 0644))
	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "ru.yml"), []byte(ruContent), 0644))

	loadTranslations()

	// Test English translations
	s.Equal("Test welcome", Get(English, "welcome"))
	s.Equal("Test check usage", Get(English, "check_usage"))

	// Test Russian translations
	s.Equal("Тестовое приветствие", Get(Russian, "welcome"))
	s.Equal("Тестовая инструкция", Get(Russian, "check_usage"))
}

func (s *TranslationsTestSuite) TestGetWithArgs() {
	// Create test YAML files
	enContent := `error_checking: "Error: %v"
result_suspicious: "Risk: %.2f, Details: %s"`

	ruContent := `error_checking: "Ошибка: %v"
result_suspicious: "Риск: %.2f, Детали: %s"`

	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "en.yml"), []byte(enContent), 0644))
	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "ru.yml"), []byte(ruContent), 0644))

	loadTranslations()

	// Test cases for error messages
	errorCases := []struct {
		lang     Language
		expected string
	}{
		{English, "Error: test error"},
		{Russian, "Ошибка: test error"},
	}

	for _, tc := range errorCases {
		s.Run(string(tc.lang), func() {
			s.Equal(tc.expected, Get(tc.lang, "error_checking", "test error"))
		})
	}

	// Test cases for result messages
	resultCases := []struct {
		lang     Language
		expected string
	}{
		{English, "Risk: 0.75, Details: suspicious activity"},
		{Russian, "Риск: 0.75, Детали: suspicious activity"},
	}

	for _, tc := range resultCases {
		s.Run(string(tc.lang), func() {
			s.Equal(tc.expected,
				Get(tc.lang, "result_suspicious", 0.75, "suspicious activity"))
		})
	}
}

func (s *TranslationsTestSuite) TestGetFallback() {
	// Create test YAML files
	enContent := `welcome: "Test welcome"
check_usage: "Test check usage"`

	ruContent := `welcome: "Тестовое приветствие"`

	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "en.yml"), []byte(enContent), 0644))
	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "ru.yml"), []byte(ruContent), 0644))

	loadTranslations()

	// Test cases for fallback behavior
	fallbackCases := []struct {
		name     string
		lang     Language
		key      string
		expected string
	}{
		{"Unknown language falls back to English", "fr", "welcome", "Test welcome"},
		{"Missing key in Russian falls back to English", Russian, "check_usage", "Test check usage"},
		{"Missing key in both languages returns empty", English, "nonexistent_key", ""},
	}

	for _, tc := range fallbackCases {
		s.Run(tc.name, func() {
			s.Equal(tc.expected, Get(tc.lang, tc.key))
		})
	}
}

func (s *TranslationsTestSuite) TestLoadTranslationsError() {
	// Create invalid YAML file
	content := `invalid: yaml: this is not valid yaml`

	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "en.yml"), []byte(content), 0644))
	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "ru.yml"), []byte(content), 0644))

	s.Panics(loadTranslations)
}

func (s *TranslationsTestSuite) TestGetWithEmptyLanguage() {
	// Create test YAML file
	content := `welcome: "Test welcome"`
	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "en.yml"), []byte(content), 0644))
	require.NoError(s.T(), os.WriteFile(filepath.Join(s.tempDir, "ru.yml"), []byte(content), 0644))

	loadTranslations()

	s.Equal("Test welcome", Get("", "welcome"))
}
