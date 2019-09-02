package main

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
)

type Cluster struct {
	Name       string   `yaml:"name"`
	DataCentre string   `yaml:"datacentre"`
	Nodes      []string `yaml:"nodes"`
}

type DoNotSend struct {
	What       string   `yaml:"what"`
	Peers string   `yaml:"peers"`
	Community       int64   `yaml:"community"`
	AS int64   `yaml:"as"`
}

type SetLocPref struct {
	Value       int   `yaml:"value"`
	Destination string   `yaml:"dest"`
	Community       int64   `yaml:"community"`
	AS int64   `yaml:"as"`
}

type NoExport struct {
	DoNotSends []DoNotSend `yaml:"donotsend"`
	SetLocPrefs []SetLocPref `yaml:"setlocpref"`
}

type SetCustomerRoute struct {
	Value       int   `yaml:"value"`
	Community       int64   `yaml:"community"`
	AS int64   `yaml:"as"`
}

type LocalPreference struct {
	SetCustomersRoute []SetCustomerRoute `yaml:"setcustomerroute"`
}

type Other struct {
	What       string   `yaml:"what"`
	Action string   `yaml:"action"`
	From string   `yaml:"from"`
	Community       int64   `yaml:"community"`
	AS int64   `yaml:"as"`
}

type Configuration struct {
	Clusters    []Cluster `yaml:"clusters"`
	Noexports    NoExport `yaml:"noexport"`
	LocalPreferences LocalPreference `yaml:"localpreference"`
	Others	[]Other `yaml:"other"`
	MinReplicas int       `yaml:"min_replicas"`
	MaxReplicas int       `yaml:"max_replicas"`
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

func main() {
    var c Configuration
    c.getConf()

    fmt.Println(c.Noexports.DoNotSends[0])
	fmt.Println(c.Noexports.DoNotSends[1])
	fmt.Println(c.Noexports.SetLocPrefs[1])
	fmt.Println(c.LocalPreferences)
	fmt.Println(c.Others)
}