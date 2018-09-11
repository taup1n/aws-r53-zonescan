package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Route53ZoneJSON struct {
	ResourceRecordSets []struct {
		AliasTarget struct {
			HostedZoneID         string `json:"HostedZoneId"`
			EvaluateTargetHealth bool   `json:"EvaluateTargetHealth"`
			DNSName              string `json:"DNSName"`
		} `json:"AliasTarget,omitempty"`
		Type            string `json:"Type"`
		Name            string `json:"Name"`
		ResourceRecords []struct {
			Value string `json:"Value"`
		} `json:"ResourceRecords,omitempty"`
		TTL           int    `json:"TTL,omitempty"`
		Weight        int    `json:"Weight,omitempty"`
		SetIdentifier string `json:"SetIdentifier,omitempty"`
		GeoLocation   struct {
			ContinentCode string `json:"ContinentCode"`
		} `json:"GeoLocation,omitempty"`
	} `json:"ResourceRecordSets"`
}

var err error

func main() {
	var route53ZoneJSON Route53ZoneJSON

	if len(os.Args[1:]) < 1 {
		fmt.Println(os.Args[0] + " [aws json file]")
		os.Exit(1)
	}

	argsWithoutProg := os.Args[1:]

	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		fmt.Println(os.Args[0] + " [aws json file]")
		os.Exit(1)
		//panic(err)
	}

	json.NewDecoder(file).Decode(&route53ZoneJSON)

	for _, record := range route53ZoneJSON.ResourceRecordSets {
		if record.Type == "A" {
			if len(record.ResourceRecords) > 0 {
				//fmt.Println(record.Name + " " + record.ResourceRecords[0].Value)

				//cmd := exec.Command("nmap", "-sS", record.Name)
				cmd := exec.Command("nmap", record.Name)
				out, _ := cmd.CombinedOutput()

				//fmt.Printf("%s\n", string(out))
				//if strings.LastIndex(string(out), "All 1000") > 0 {
				if strings.LastIndex(string(out), "Host seems down") > 0 {
					fmt.Println(record.Name + " " + record.ResourceRecords[0].Value)
				}
			}
		}
	}
}
