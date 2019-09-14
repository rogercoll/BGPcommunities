package parserNaturalLang

import (
	"log"
	"context"
	"testing"
	language "cloud.google.com/go/language/apiv1"
)

func TestAnalizeSyntax(t *testing.T){
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
			log.Fatalf("Failed to create client: %v", err)
	}
	// Sets the text to analyze.
	text := readFromFile("examples/as174.txt")
	m, err := analyzeSyntax(ctx,client,text)
	if err != nil {
		log.Fatal(err)
	}
	//printSentences(m)
	parserCommunities(m)
}


