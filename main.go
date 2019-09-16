package main

import (
	"os"
	"log"
	"context"
	"io/ioutil"
    "gopkg.in/yaml.v2"
	"github.com/rogercoll/BGPcommunities/parserNaturalLang"
)

type DoNotSend struct {
	What       string   `yaml:"what"`
	Peers string   `yaml:"peers"`
	Community       int64   `yaml:"community"`
}

type SetLocPref struct {
	Value       int   `yaml:"value"`
	Destination string   `yaml:"dest"`
	Community       int64   `yaml:"community"`
}

type NoExport struct {
	DoNotSends []DoNotSend `yaml:"donotsend"`
	SetLocPrefs []SetLocPref `yaml:"setlocpref"`
}

type SetCustomerRoute struct {
	Value       int   `yaml:"value"`
	Community       int64   `yaml:"community"`
}

type LocalPreference struct {
	SetCustomersRoute []SetCustomerRoute `yaml:"setcustomerroute"`
}

type DoNotAnnounce struct {
	Peer string `yaml:"peer"`
	Asn int64	`yaml:"asn"`
	Community int64 `yaml:"community"`
}

type Prepend struct {
	What string	`yaml:"what`
	Times int	`yaml:"times"`
	Community int64 `yaml:"community"`
}

type PeerControl struct {
	DoNotAnnounces []DoNotAnnounce `yaml:"donotannounce"`
	Prepends []Prepend `yaml:"prepend"`
}

type Other struct {
	What       string   `yaml:"what"`
	Action string   `yaml:"action"`
	From string   `yaml:"from"`
	Community       int64   `yaml:"community"`
}

type Configuration struct {	
	Noexports    NoExport `yaml:"noexport"`
	LocalPreferences LocalPreference `yaml:"localpreference"`
	PeerControls PeerControl `yaml:"peercontrol"`
	Others	[]Other `yaml:"other"`
	As int64       `yaml:"as"`
}

func (c *Configuration) getConf() *Configuration {

    yamlFile, err := ioutil.ReadFile("conf.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
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

func main() {
    var c Configuration
	c.getConf()
	ctx := context.Background()
	client, err := parserNaturalLang.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	text := readFromFile("parserNaturalLang/examples/as174.txt")
	m, err := parserNaturalLang.AnalyzeSyntax(ctx,client,text)
	if err != nil {
		log.Fatal(err)
	}
	parserNaturalLang.ParserCommunities(m)
}
