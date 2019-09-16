package parserNaturalLang

import (
	"os"
	"fmt"
	"log"
	"context"
	"strconv"
	"github.com/golang/protobuf/proto"
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"github.com/rogercoll/BGPcommunities/communities"
)

var noExports communities.NoExport
var localPreference	communities.LocalPreference
var peerControls	communities.PeerControl
var other	communities.Other

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
	//fmt.Println(len(sentence))
	var err error
	var as,community int
	foundDoublePoint := false
	for _,token := range sentence {
		fmt.Printf("%v ", token.GetLemma())
		tTag := token.GetPartOfSpeech().GetTag()
		//11 equals to VERB
		if tTag == 11 {
			fmt.Printf("Verb found: %s\n", token.GetLemma())
		}
		if tTag == 10 {
			if token.GetLemma() == ":" {
				foundDoublePoint = true
			}
		}
		//7 equals to NUM
		if tTag == 7 {
			if foundDoublePoint {
				community, err = strconv.Atoi(token.GetLemma())
				if err != nil {
					log.Fatal(err)
				}
			} else {
				as, err = strconv.Atoi(token.GetLemma())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	fmt.Println()
	fmt.Println(as)
	fmt.Println(community)
}


func ParserCommunities(m *languagepb.AnnotateTextResponse) {
	//conf := communities.NoExport{}
	tokens := m.GetTokens()
	var tokensForSentence []*languagepb.Token
	for _,token := range tokens {
		tTag := token.GetPartOfSpeech().GetTag()
		//10 equals to PUNCT
		if tTag == 10 && token.GetLemma() == "." {
			//THIS CAN BE DONE IN PARALLEL
			parseSentenceTokens(tokensForSentence)
			tokensForSentence = nil
		} else {
			tokensForSentence = append(tokensForSentence, token)
			continue
		}
		
	}
}


func AnalyzeSyntax(ctx context.Context, client *language.Client, text string) (*languagepb.AnnotateTextResponse, error) {
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

func NewClient(ctx context.Context) (*language.Client, error) {
	return language.NewClient(ctx)
}