package mocks

import (
	"context"
	"encoding/csv"
	"errors"
	"golang.org/x/text/language"
	"os"
)

type translationEntry struct {
	SourceLanguage language.Tag
	SourceText     string
	TargetLanguage language.Tag
	TranslatedText string
}

type MockTranslator struct {
	// mock-specific fields
	translations map[string]translationEntry
}

func NewMockTranslator(csvFilePath string) (*MockTranslator, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	mt := &MockTranslator{
		translations: make(map[string]translationEntry),
	}

	for _, record := range records[1:] {
		srcLang, _ := language.Parse(record[0])
		tgtLang, _ := language.Parse(record[2])
		mt.translations[record[1]+record[0]+record[2]] = translationEntry{
			SourceLanguage: srcLang,
			SourceText:     record[1],
			TargetLanguage: tgtLang,
			TranslatedText: record[3],
		}
	}

	return mt, nil
}

func (m *MockTranslator) Translate(ctx context.Context, texts []string, targetLang language.Tag) ([]string, error) {
	var translations []string
	for _, text := range texts {
		key := text + "en" + targetLang.String()
		if translation, ok := m.translations[key]; ok {
			translations = append(translations, translation.TranslatedText)
		} else {
			return nil, errors.New("translation not found")
		}
	}
	return translations, nil
}
