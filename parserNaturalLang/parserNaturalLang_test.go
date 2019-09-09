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
	text := "NO_EXPORT.\nDo not send route to NA 174:970.\n"
	m, err := analyzeSyntax(ctx,client,text)
	if err != nil {
		log.Fatal(err)
	}
	printSentences(m)
}


