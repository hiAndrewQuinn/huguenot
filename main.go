package main

import (
	"context"
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

func printAST(node ast.Node, level int) {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}
	log.Printf("%sNode Type: %T\n", indent, node)
	if node.HasChildren() {
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			printAST(child, level+1)
		}
	}
}

func translateAST(ctx context.Context, client *translate.Client, node ast.Node, indent string, source []byte) {
	log.Printf("%sNode Type: %T", indent, node)

	language := language.English

	if textNode, ok := node.(*ast.Text); ok {
		segment := textNode.Segment
		text := string(source[segment.Start:segment.Stop])
		log.Printf("%sText: %s\n", indent, text)
		// translate the text and print that too

		// Perform translation and handle the response
		translations, err := client.Translate(ctx, []string{text}, language, nil)
		if err != nil {
			log.Fatalf("Failed to translate text: %v", err)
		}
		if len(translations) > 0 {
			log.Printf("%sTrns: %s\n", indent, translations[0].Text)
		}
	}

	if node.HasChildren() {
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			translateAST(ctx, client, child, indent+"  ", source)
		}
	}
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
	translateAST(ctx, client, document, "", mdContent)
	log.Printf("---")
}
