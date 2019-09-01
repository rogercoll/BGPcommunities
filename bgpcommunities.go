package main

import(
	"fmt"
	"gopkg.in/yaml.v2"
    	"io/ioutil"
        "os"
)

type SetLocalPreference struct {
	Value int
	Dest string
	Community int
	As int
}

type DoNotSend struct {  
	What string
	Peers string
	Community int
	As int
}


type DoNotAnnounce struct {
	Peers string
}

type NoExport struct{
	Sets []SetLocalPreference
	NoSend []DoNotSend 
}

type SetCustomerRoute struct{
	Value int
	Community int
	As int
}

type LocPref struct{
	Sets []SetCustomerRoute 
}

type PeerContr struct {
	NoAnns []DoNotAnnounce
}

type Other struct {
	What string
	Action string
	Peer string
	Community int
	As int
}



type Config struct{
	NoexpComms []NoExport
	LocPrefComms []LocPref
	PeerComms []PeerContr
	OtherComms []Other
} 

func Hello()string {
	fmt.Println("hey")
	return "hello"
}


func main(){
	filename := os.Args[1]
	 var config Config
	source, err := ioutil.ReadFile(filename)
	if err != nil {
			            panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {       panic(err)
	}
	fmt.Printf("Value: %#v\n", config.OtherComms[0].As)

} 
