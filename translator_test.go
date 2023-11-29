package main

import (
	"context"
	"golang.org/x/text/language"
	"huguenot/mocks"
	"testing"
)

func TestMockTranslator(t *testing.T) {
	// Create the mock translator
	mockTranslator, err := mocks.NewMockTranslator("mocks/mock_translations.csv")
	if err != nil {
		t.Fatalf("Failed to create mock translator: %v", err)
	}

	// Define a test case
	tests := []struct {
		name           string
		sourceText     string
		targetLang language.Tag
		want           string
	}{
		{
			name:           "English to Spanish",
			sourceText:     "Hello",
			targetLang: language.Spanish,
			want:           "Hola",
		},
		// Add more test cases as needed
	}

	// Run the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := mockTranslator.Translate(context.Background(), []string{tc.sourceText}, tc.targetLang)
			if err != nil {
				t.Errorf("Translate() error = %v", err)
				return
			}
			if got[0] != tc.want {
				t.Errorf("Translate() got = %v, want %v", got[0], tc.want)
			}
		})
	}
}

