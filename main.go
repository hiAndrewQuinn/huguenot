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

func checkArgs(args []string) string {
	if len(args) != 2 {
		log.Fatalf("Usage: %s <markdown-file>", args[0])
	}
	return args[1]
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

func translateAST(ctx context.Context, client *translate.Client, language language.Tag, node ast.Node, indent string, source []byte) error {
	log.Printf("%sNode Type: %T", indent, node)

	if textNode, ok := node.(*ast.Text); ok {
		segment := textNode.Segment
		originalText := string(source[segment.Start:segment.Stop])
		log.Printf("%sText: %s\n", indent, originalText)

		// Perform translation
		translations, err := client.Translate(ctx, []string{originalText}, language, nil)
		if err != nil {
			return fmt.Errorf("failed to translate text: %v", err)
		}
		if len(translations) > 0 {
			translatedText := translations[0].Text

			// Create a new text node with the translated text
			newNode := ast.NewTextSegment(text.NewSegment(segment.Start, segment.Start+len(translatedText)))
			newNode.AppendChild(newNode, ast.NewString([]byte(translatedText)))

			// Replace the original text node with the new one
			if parent := textNode.Parent(); parent != nil {
				parent.ReplaceChild(parent, textNode, newNode)
			}
		}
	}

	if node.HasChildren() {
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			err := translateAST(ctx, client, language, child, indent+"  ", source)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	ctx := context.Background()

	client := makeClient(ctx, "service-account-key.json")
	defer client.Close()

	mdContent := getMdContent(checkArgs(os.Args))

	markdown := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	document := markdown.Parser().Parse(text.NewReader(mdContent))
	err := translateAST(ctx, client, language.English, document, "", mdContent)
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
