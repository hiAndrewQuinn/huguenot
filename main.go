package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/translate"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

func checkArgs(args []string) (string, language.Tag) {
	if len(args) != 3 {
		log.Fatalf("Usage: %s <markdown-file> <language-code>", args[0])
	}
	mdFile := args[1]
	langCode := args[2]

	// Parse the language code
	lang, err := language.Parse(langCode)
	if err != nil {
		log.Fatalf("Invalid language code: %v", err)
	}

	return mdFile, lang
}

func makeClient(ctx context.Context, service_account_keyfile string) *translate.Client {
	client, err := translate.NewClient(ctx, option.WithCredentialsFile(service_account_keyfile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func getMdContent(markdownFile string) []byte {
	mdContent, err := os.ReadFile(markdownFile)
	if err != nil {
		log.Fatalf("Failed to read markdown file: %v", err)
	}
	return mdContent
}

func getStringsToTranslate(node ast.Node, source []byte) []string {
	// what is the golang equivalent of Python's pass?
	var stringsToTranslate []string

	if textNode, ok := node.(*ast.Text); ok {
		segment := textNode.Segment
		originalText := string(source[segment.Start:segment.Stop])
		stringsToTranslate = append(stringsToTranslate, originalText)
	} else if node.HasChildren() {
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			// Shoot, now that's a good question. Does Go have pass by value, or pass by reference?
			// I guess pass by value, since we have pointers and stuff here.
			childStrings := getStringsToTranslate(child, source)
			stringsToTranslate = append(stringsToTranslate, childStrings...)
		}
	}
	return stringsToTranslate
}

func translateAST(ctx context.Context, client *translate.Client, language language.Tag, node ast.Node, indent string, source []byte) error {
	stringsToTranslate := getStringsToTranslate(node, source)
	if len(stringsToTranslate) == 0 {
		log.Print("No strings to translate found in this document.")
	} else {
		log.Printf("List of strings to translate...")
		for i := 0; i < len(stringsToTranslate); i++ {
			log.Printf("- %s", stringsToTranslate[i])
		}
	}

	// Let's send the whole stringsToTranslate up as a single translation call.
	translation, err := client.Translate(ctx, stringsToTranslate, language, nil)
	if err != nil {
		return fmt.Errorf("failed to translate text: %v", err)
	}

	// And print it out.
	log.Printf("Translated text: %s", translation)

	////		if len(translations) > 0 {
	////			translatedText := translations[0].Text
	////
	////			// Create a new text node with the translated text
	////			newNode := ast.NewTextSegment(text.NewSegment(segment.Start, segment.Start+len(translatedText)))
	////			newNode.AppendChild(newNode, ast.NewString([]byte(translatedText)))
	////
	////			// Replace the original text node with the new one
	////			if parent := textNode.Parent(); parent != nil {
	////				parent.ReplaceChild(parent, textNode, newNode)
	////			}
	////		}
	//	}
	//
	//	if node.HasChildren() {
	//		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
	//			err := translateAST(ctx, client, language, child, indent+"  ", source)
	//			if err != nil {
	//				return err
	//			}
	//		}
	//	}
	return nil
}

func main() {
	ctx := context.Background()

	// print a Hello! in yellow text.
	log.Default().Printf("\033[33mHello!\033[0m")

	client := makeClient(ctx, "service-account-key.json")
	defer client.Close()

	// Print Google Tranlsate's supported languages.
	// func (c *Client) SupportedLanguages(ctx context.Context, target language.Tag) ([]Language, error)
	log.Default().Printf("\033[33mSupported Languages:\033[0m")

	langs, err := client.SupportedLanguages(ctx, language.English)
	if err != nil {
		log.Fatalf("Failed to get supported languages: %v", err)
	}

	for i := 0; i < len(langs); i++ {
		log.Default().Printf("\033[33m%s\033[0m", langs[i].Tag)
	}

	mdFile, lang := checkArgs(os.Args)
	mdContent := getMdContent(mdFile)

	markdown := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	document := markdown.Parser().Parse(text.NewReader(mdContent))
	err = translateAST(ctx, client, lang, document, "", mdContent)
	if err != nil {
		log.Fatalf("Error during translation: %v", err)
	}

	var buf bytes.Buffer
	if err := markdown.Renderer().Render(&buf, mdContent, document); err != nil {
		log.Fatalf("Failed to render markdown: %v", err)
	}

	err = os.WriteFile("output.md", buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	log.Printf("---")
}
