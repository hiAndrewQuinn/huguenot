package main

import (
	"context"
	"golang.org/x/text/language"
)

type translator interface {
	Translate(ctx context.Context, texts []string, targetLang language.Tag) ([]string, error)
}
