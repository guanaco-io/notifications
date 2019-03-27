package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var config = &Configuration

type Alert struct {
	Id string `json:"id"`
	Text string  `json:"text"`
	Resource string `json:"resource"`
	Severity string `json:"severity"`
	Event string `json:"event"`
	Attributes map[string]string `json:"attributes"`
}

type AlertsResponse struct {
	Alerts []Alert `json:"alerts"`
	StatusCounts map[string]int `json:"statusCounts"`
}

func countOpenAlerts() int {
	url := fmt.Sprintf("%v/alerts", config.Alerta.Endpoint)
	resp, err := performRequest("GET", url, nil)

	if err != nil {
		log.Printf("Error fetching alert count: %v", err)
		return 0
	}

	log.Printf("Response: %s", resp.Status)

	decoder := json.NewDecoder(resp.Body)

	var alertsResponse AlertsResponse
	if err := decoder.Decode(&alertsResponse); err != nil {
		log.Print(err)
		log.Printf("Error parsing alerts response: %v", err)
		return 0
	}

	closeError := resp.Body.Close()
	if closeError != nil {
		log.Fatalf("Error closing response body: %v", closeError)
	}

	return len(alertsResponse.Alerts)
}

func performRequest(method string, url string, body []byte) (resp *http.Response, err error) {
	log.Printf("Doing %s request to %s with body: %s", method, url, body)

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}
