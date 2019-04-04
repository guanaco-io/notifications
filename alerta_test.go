package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestUnMarshallAlertsResponse(t *testing.T) {

	data, readFileError := ioutil.ReadFile("test/alerts.json")
	if readFileError != nil {
		t.Fatalf("cannot read test file test/alerts.json: %v", readFileError)
	}

	var alertsResponse AlertsResponse

	err := json.Unmarshal(data, &alertsResponse)

	if err != nil {
		t.Fatalf("cannot unmarshal data: %v", err)
	}

	fmt.Println(alertsResponse)
	log.Printf("Parsed response: %v", alertsResponse)
}