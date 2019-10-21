package parserNaturalLang

import (
	"os"
	"fmt"
	"log"
	"context"
	"strconv"
	"gopkg.in/yaml.v2"
	"github.com/golang/protobuf/proto"
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"github.com/rogercoll/BGPcommunities/communities"
)

type Configuration struct {	
	Noexports    communities.NoExport `yaml:"noexport,omitempty"`
	LocalPreferences communities.LocalPreference `yaml:"localpreference,omitempty"`
	PeerControls communities.PeerControl `yaml:"peercontrol,omitempty"`
	Others	[]communities.Other `yaml:"other,omitempty"`
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
			if token.GetText().GetContent() == "to" {
				if doNotSent {
					nextPeer = true
				} else {
					nextLocValue = true
				}
			} else if token.GetText().GetContent() == "in" {
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
		if token.GetText().GetContent() == "tp" || token.GetText().GetContent() == "pref"{
			nextLocValue = true
		}
		if tTag == 2 {
			if token.GetText().GetContent() == "to" {
				nextLocValue = true
			} else if token.GetText().GetContent() == "in" {
				nextPeer = true
			}
		}
		//6 equals to Noun
		if tTag == 6 {
			if nextPeer {
				peers += token.GetLemma() + " "
			}
		}
		if token.GetText().GetContent() == "="{
			nextLocValue = true
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
	var err error
	var as,community,asnValue int
	nextPeer, prepend := false, false
	peers := ""
	foundDoublePoint := false
	for _,token := range sentence {
		tTag := token.GetPartOfSpeech().GetTag()
		fmt.Println(token.GetText().GetContent())
		if token.GetText().GetContent() == "Prepend" || token.GetText().GetContent() == "prepend" {
			prepend = true
		}
		if tTag == 11 {
			if token.GetLemma() == "Prepend" || token.GetText().GetContent() == "Prepend" {
				prepend = true
			}
		}
		//2 equals to Adposition (preposition and postposition)
		if tTag == 2 {
			if token.GetLemma() == "to" || token.GetLemma() == "in" {
				nextPeer = true
				continue
			}
		}
		if tTag == 10 {
			if token.GetLemma() == ":" {
				foundDoublePoint = true
			}
		}
		//7 equals to NUM
		if tTag == 7 {
			if nextPeer{
				asnValue, err = strconv.Atoi(token.GetLemma())
				if err != nil {
					return err
				}
				nextPeer = false
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
		aux := token.GetText().GetContent()
			if (prepend && (string(aux[len(aux)-1]) == "x")) {
				asnValue, err = strconv.Atoi(aux[:len(aux)-1])
				if err != nil {
					return err
				}
				continue
			}
		if nextPeer {
			peers += token.GetLemma() + " "
			continue
		}
	}
	fmt.Println()
	newConf.As = as
	if prepend {
		comm := communities.Prepend{What: peers, Times: asnValue, Community: community}
		newConf.PeerControls.Prepends = append(newConf.PeerControls.Prepends, comm)
	} else {
		comm := communities.DoNotAnnounce{Peer: peers, Asn: asnValue, Community: community}
		newConf.PeerControls.DoNotAnnounces = append(newConf.PeerControls.DoNotAnnounces, comm)
		fmt.Println(comm)
	}
	return nil
}

func parseSentencesOther(sentence []*languagepb.Token) (error) {
	var err error
	from, what, verbs := "", "", ""
	fromdetect, foundDoublePoint := false, false
	var community int
	for _,token := range sentence {
		tTag := token.GetPartOfSpeech().GetTag()
		if fromdetect {
			from = token.GetLemma()
			fromdetect = false
			continue
		}
		if tTag == 2 && token.GetLemma() == "from" {
			fromdetect = true
		}
		if tTag == 10 {
			if token.GetLemma() == ":" {
				foundDoublePoint = true
			}
		}
		if tTag == 7 {
			if foundDoublePoint {
				community, err = strconv.Atoi(token.GetLemma())
				if err != nil {
					return err
				}
			}
		}
		if tTag == 6 || tTag == 1 {
			what += token.GetText().GetContent()
		}
		if tTag == 11 {
			verbs += token.GetText().GetContent()
		}
	}
	comm := communities.Other{From: from, Community: community, What: what, Action: verbs}
	newConf.Others = append(newConf.Others, comm)
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
	newConf = Configuration{}
	tokens := m.GetTokens()
	var tokensForSentence []*languagepb.Token
	action := ""
	for _,token := range tokens {
		tTag := token.GetPartOfSpeech().GetTag()
		//10 equals to PUNCT
		if token.GetLemma() == "NO_EXPORT" {
			action = token.GetLemma()
		} else if token.GetLemma() == "LOCAL_PREFERENCE" || token.GetLemma() == "LOCAL_PREFERENCE." {
			action = token.GetLemma()
		} else if token.GetLemma() == "PEER_CONTROLS" {
			action = token.GetLemma()
		} else if token.GetLemma() == "OTHER_COMMUNITIES" || token.GetLemma() == "OTHER_COMMUNITIES." {
			action = token.GetLemma()
		} else {
			if tTag == 10 && token.GetLemma() == "." {
				//THIS CAN BE DONE IN PARALLEL
				if action == "NO_EXPORT" {
					parseSentencesNoExport(tokensForSentence)
					tokensForSentence = nil
				} else if action == "LOCAL_PREFERENCE" || action == "LOCAL_PREFERENCE." {
					parseSentencesLocalPreference(tokensForSentence)
					tokensForSentence = nil
				} else if action == "PEER_CONTROLS" {
					parseSentencesPeerControls(tokensForSentence)
					tokensForSentence = nil
				} else if action == "OTHER_COMMUNITIES" || action == "OTHER_COMMUNITIES." {
					parseSentencesOther(tokensForSentence)
					tokensForSentence = nil
				}
				
			} else {
				tokensForSentence = append(tokensForSentence, token)
				continue
			}
			
		}
	}
	b, err := yaml.Marshal(newConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	f, err := os.Create("yaml/as" + strconv.Itoa(newConf.As))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
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