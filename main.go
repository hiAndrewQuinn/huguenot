package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Use the service account JSON key file
	client, err := translate.NewClient(ctx, option.WithCredentialsFile("./service-account-key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	text := "Hello, World!"
	// targetLang := "es" // Spanish language code

	translations, err := client.Translate(ctx, []string{text}, language.Spanish, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	for _, translation := range translations {
		fmt.Printf("Translated text: %v\n", translation.Text)
	}
}
