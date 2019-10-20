package main

import (
	"os"
	"log"
	"context"
	"io/ioutil"
    //"gopkg.in/yaml.v2"
	"github.com/rogercoll/BGPcommunities/parserNaturalLang"
)

/*

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


*/

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
	/*
    var c Configuration
	c.getConf()
	*/
	dir := "/home/neck/Documents/Uni/TFG/BGPcommunities/parserNaturalLang/examples/"
	files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
	}
	for _,f := range files {
		log.Println(f.Name())
		text := ""
		text = readFromFile(dir + f.Name())
		log.Println(text)
		ctx := context.Background()
		client, err := parserNaturalLang.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
		}
		m, err := parserNaturalLang.AnalyzeSyntax(ctx,client,text)
		if err != nil {
			log.Fatal(err)
		}
		parserNaturalLang.ParserCommunities(m)
	}
	
}
