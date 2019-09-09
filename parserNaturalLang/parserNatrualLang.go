package parserNaturalLang

import (
	"os"
	"log"
	"context"
	"github.com/golang/protobuf/proto"
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
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

