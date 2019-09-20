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

type Configuration struct {	
	Noexports    communities.NoExport `yaml:"noexport"`
	LocalPreferences communities.LocalPreference `yaml:"localpreference"`
	PeerControls communities.PeerControl `yaml:"peercontrol"`
	Others	[]communities.Other `yaml:"other"`
	As int       `yaml:"as"`
}

var newConf Configuration

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

func parseSentencesNoExport(sentence []*languagepb.Token) (error) {
	//comm := communities.NoExport{}
	fmt.Println("Hey NoExport sentence")
	var err error
	var as,community,locPrefValue int
	doNotSent, setLocPref, nextPeer, nextLocValue := false, false, false, false
	peers := ""
	foundDoublePoint := false
	for _,token := range sentence {
		//fmt.Printf("Lemma: %v ", token.GetLemma())
		tTag := token.GetPartOfSpeech().GetTag()
		//fmt.Printf("Tag :%v ", tTag)
		//11 equals to VERB
		if tTag == 11 {
			if token.GetLemma() == "Do" {
				doNotSent = true
			} else if token.GetLemma() == "Set" {
				setLocPref = true
			}
		}
		//2 equals to Adposition (preposition and postposition)
		if tTag == 2 {
			if token.GetLemma() == "to" {
				if doNotSent {
					nextPeer = true
				} else {
					nextLocValue = true
				}
			} else if token.GetLemma() == "in" {
				nextPeer = true
			}
		}
		//6 equals to Noun
		if tTag == 6 {
			if nextPeer {
				peers += token.GetLemma() + " "
			}
		}
		if tTag == 10 {
			if token.GetLemma() == ":" {
				foundDoublePoint = true
			}
		}
		//7 equals to NUM
		if tTag == 7 {
			if nextLocValue {
				locPrefValue, err = strconv.Atoi(token.GetLemma())
				if err != nil {
					return err
				}
				nextLocValue = false
			} else {
				if foundDoublePoint {
					community, err = strconv.Atoi(token.GetLemma())
					if err != nil {
						return err
					}
				} else {
					as, err = strconv.Atoi(token.GetLemma())
					if err != nil {
						return err
					}
				}
			}
		}
	}
	fmt.Println()
	newConf.As = as
	if setLocPref {
		comm := communities.SetLocPref{Value: locPrefValue, Community: community, Destination: peers}
		newConf.Noexports.SetLocPrefs = append(newConf.Noexports.SetLocPrefs, comm)
		fmt.Println(comm)
	} else {
		comm := communities.DoNotSend{What: "route", Peers: peers, Community: community}
		newConf.Noexports.DoNotSends = append(newConf.Noexports.DoNotSends, comm)
		fmt.Println(comm)
	}
	return nil
}

func parseSentencesLocalPreference(sentence []*languagepb.Token) (error) {
	//comm := communities.LocalPreference{}
	fmt.Println("Hey LocalPreference sentence")
	var err error
	var as,community,locPrefValue int
	nextPeer, nextLocValue := false, false
	peers := ""
	foundDoublePoint := false
	for _,token := range sentence {
		//fmt.Printf("Lemma: %v ", token.GetLemma())
		tTag := token.GetPartOfSpeech().GetTag()
		//fmt.Printf("Tag :%v ", tTag)
		//2 equals to Adposition (preposition and postposition)
		if tTag == 2 {
			if token.GetLemma() == "to" {
				nextLocValue = true
			} else if token.GetLemma() == "in" {
				nextPeer = true
			}
		}
		//6 equals to Noun
		if tTag == 6 {
			if nextPeer {
				peers += token.GetLemma() + " "
			}
		}
		if tTag == 10 {
			if token.GetLemma() == ":" {
				foundDoublePoint = true
			}
		}
		//7 equals to NUM
		if tTag == 7 {
			if nextLocValue {
				locPrefValue, err = strconv.Atoi(token.GetLemma())
				if err != nil {
					return err
				}
				nextLocValue = false
			} else {
				if foundDoublePoint {
					community, err = strconv.Atoi(token.GetLemma())
					if err != nil {
						return err
					}
				} else {
					as, err = strconv.Atoi(token.GetLemma())
					if err != nil {
						return err
					}
				}
			}
		}
	}
	fmt.Println()
	newConf.As = as
	comm := communities.SetCustomerRoute{Value: locPrefValue, Community: community}
	newConf.LocalPreferences.SetCustomersRoute = append(newConf.LocalPreferences.SetCustomersRoute, comm)
	fmt.Println(comm)
	return nil
}

func parseSentencesPeerControls(sentence []*languagepb.Token) (error) {
	//comm := communities.PeerControl{}
	fmt.Println("Hey PeerControls sentence")

	return nil
}

func parseSentencesOther(sentence []*languagepb.Token) (error) {
	//comm := communities.Other{}
	fmt.Println("Hey Other sentence")

	return nil
}


func parseSentenceTokens(sentence []*languagepb.Token)  {
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
	action := ""
	for _,token := range tokens {
		tTag := token.GetPartOfSpeech().GetTag()
		//10 equals to PUNCT
		if token.GetLemma() == "NO_EXPORT" {
			action = token.GetLemma()
		} else if token.GetLemma() == "LOCAL_PREFERENCE" {
			action = token.GetLemma()
		} else if token.GetLemma() == "PEER_CONTROLS" {
			action = token.GetLemma()
		} else if token.GetLemma() == "OTHER_COMMUNITIES" {
			action = token.GetLemma()
		} else {
			if tTag == 10 && token.GetLemma() == "." {
				//THIS CAN BE DONE IN PARALLEL
				if action == "NO_EXPORT" {
					parseSentencesNoExport(tokensForSentence)
					tokensForSentence = nil
				} else if action == "LOCAL_PREFERENCE" {
					parseSentencesLocalPreference(tokensForSentence)
					tokensForSentence = nil
				} else if action == "PEER_CONTROLS" {
					parseSentencesPeerControls(tokensForSentence)
					tokensForSentence = nil
				} else if action == "OTHER_COMMUNITIES" {
					parseSentencesOther(tokensForSentence)
					tokensForSentence = nil
				}
				
			} else {
				tokensForSentence = append(tokensForSentence, token)
				continue
			}
			
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