package parserNaturalLang

import (
	"os"
	"fmt"
	"log"
	"context"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	//"github.com/rogercoll/BGPcommunities/communities"
)

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}

func storeToFile(v proto.Message, err error) {
	f, err := os.Create("/home/rcoll/freeFridays/BGPcommunities/parserNaturalLang/info.txt")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(f, v)
}

func readFromFile(path string) string {
	file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
  	b, err := ioutil.ReadAll(file)
	return string(b)
}

func printSentences(m *languagepb.AnnotateTextResponse) {
	sentences := m.GetSentences()
	for _,sentence := range sentences {
		content := sentence.GetText()
		if content.GetContent() == "NO_EXPORT." {
			fmt.Println("Heey no export sentences found ")
		}
	}
}

func parseSentenceTokens(sentence []*languagepb.Token) {
	for token := range sentence {

	}
}


func parserCommunities(m *languagepb.AnnotateTextResponse) {
	//conf := communities.NoExport{}
	tokens := m.GetTokens()
	totalSize := len(tokens)
	var tokensForSentence []*languagepb.Token
	for i,token := range tokens {
		tTag := token.GetPartOfSpeech().GetTag()
		//10 equals to PUNCT
		if tTag != 10 && token.GetLemma() != "." {
			tokensForSentence = append(tokensForSentence, token)
			continue
		} else {
			parseSentenceTokens(tokensForSentence)
			tokensForSentence = nil
		}
		//11 equals to VERB
		if tTag == 11 {
			fmt.Printf("Verb found: %s\n", token.GetLemma())
			if ((i + 3) < totalSize) {
				i += 3
			}
		}
	}
}


func analyzeSyntax(ctx context.Context, client *language.Client, text string) (*languagepb.AnnotateTextResponse, error) {
	return client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

