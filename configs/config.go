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
	RevenueViaCustomVariable      struct {
		Key                string `json:"Key"`
		CustomVariableName string `json:"CustomVariableName"`
		TrafficSourceId    string `json:"TrafficSourceId"`
		FieldName          string `json:"FieldName"`
		APIVariableName    string `json:"ApiVariableName"`
	} `json:"RevenueViaCustomVariable"`
	TSMappingViaCustomVariable struct {
		Key                string `json:"Key"`
		CustomVariableName string `json:"CustomVariableName"`
		TrafficSourceId    string `json:"TrafficSourceId"`
		FieldName          string `json:"FieldName"`
		APIVariableName    string `json:"ApiVariableName"`
	} `json:"TSMappingViaCustomVariable"`
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
		fmt.Println("Follwing values has been picked from config values:")
		fmt.Println(Configurations)
	}
}
