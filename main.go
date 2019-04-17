package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	maxDuration time.Duration = 1<<63 - 1
)

func main() {

	fmt.Printf("Starting notifications client")
	if len(os.Args) < 2 {
		log.Printf("Usage: notifications <config.yml>")
		log.Fatal("  <config.yml> parameter is missing!")
	}

	config, initError := Load(os.Args[1])
	logFatal("Error initializing program", initError)
	fmt.Printf("Configuration loaded successfully")

	channels, channelsError := LoadChannels(config)
	logFatal("Error loading channels configuration", channelsError)
	fmt.Printf("%v Channels loaded successfully", len(channels))


	client := AlertaClient{config: config.Alerta}
	ticker := time.NewTicker(config.Alerta.ReloadInterval * time.Second)

	fmt.Printf("Waiting for %v before fetching alerts", config.Alerta.ReloadInterval * time.Second)
	go func() {
		for t := range ticker.C {

			log.Printf("Fetching Alerta alerts at %v", t)

			for ruleName, rule := range config.Rules {

				log.Printf("Evaluating rule %v", ruleName)

				if alerts := client.searchAlerts(rule); alerts != nil && len(alerts) > 0 {

					alreadyNotified, notNotified := Partition(alerts, ruleName, IsNotified)

					if nrOfAlerts := len(notNotified); nrOfAlerts > 0 {

						for _, ruleChannel := range rule.Channels {
							log.Printf("Sending %v alerts to channel %v of rule %v", nrOfAlerts, ruleChannel, ruleName)

							channel, ok := channels[ruleChannel]
							if !ok {
								log.Fatalf("Unable to find channel '%v' of rule '%v' in channel config", ruleChannel, ruleName)
							}

							sendError := channel.Send(AlertEvent{NewAlertCount: nrOfAlerts, NewAlerts: notNotified, AlreadyNotified: len(alreadyNotified)}, config.DryRun)
							if sendError != nil {
								log.Printf("Error sending alert event to channel '%v' of rule '%v': %v", ruleChannel, ruleName, sendError)
							}
						}

					}

					log.Printf("%v alerts were already notified for rule %v", len(alreadyNotified), ruleName)

				} else {
					log.Printf("No Alerts found for rule %v", ruleName)
				}

			}
		}
	}()

	time.Sleep(maxDuration) // TODO find a better way to block here
	ticker.Stop()
	log.Println("Ticker stopped, exiting program")

	os.Exit(0)
}

func Partition(all []Alert, ruleId string, predicate func(Alert, string) bool) ([]Alert, []Alert) {
	success := make([]Alert, 0)
	failure := make([]Alert, 0)
	for _, alert := range all {
		if predicate(alert, ruleId) {
			success = append(success, alert)
		} else {
			failure = append(failure, alert)
		}
	}
	return success, failure
}

// Log error message and exit program
func logFatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%v: %v", msg, err)
	}
}
