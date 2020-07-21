package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Configs struct {
	SpreadsheetId                 string
	TrafficSourcesShortlisted     []string
	TrafficSourceFilteringEnabled bool
	IncludeTrafficSources         string
	VoluumAccessId                string
	VoluumAccessKey               string
}

var (
	Configurations = Configs{}
)

func SetConfig() {
	input, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	error := json.Unmarshal(input, &Configurations)
	if error != nil {
		fmt.Println("Config file is missing in root directory")
		panic(error)
	} else {
		fmt.Println("Follwing values has been picked from ")
		fmt.Println(Configurations)
	}
}
