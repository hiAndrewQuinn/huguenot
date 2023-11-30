package main

import (
	"context"
	"golang.org/x/text/language"
)

type Lang struct {
	langCode string
	langName string
}

type translator interface {
	getDstLangs() []Lang
	Translate(ctx context.Context, texts []string, targetLang language.Tag) ([]string, error)
}
