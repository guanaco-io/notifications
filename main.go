package main

import (
	"log"
	"os"
	"time"
)

const (
	maxDuration time.Duration = 1<<63 - 1
)

func main() {

	if len(os.Args) != 2 {
		log.Printf("Usage: notifications <config.yml>")
		log.Fatal("  <config.yml> parameter is missing!")
	}
	log.Printf("Starting Guanaco notifications app with config file %v", os.Args[1])

	config, initError := Load(os.Args[1])
	logFatal("Error initializing program", initError)
	log.Printf("Configuration loaded successfully")

	channels, channelsError := LoadChannels(config)
	logFatal("Error loading channels configuration", channelsError)
	log.Printf("%v Channels loaded successfully", len(channels))

	client := AlertaClient{config: config.Alerta}
	ticker := time.NewTicker(config.Alerta.ReloadInterval * time.Second)

	log.Printf("Waiting for %v before fetching alerts", config.Alerta.ReloadInterval*time.Second)

	ruleHandlers := make([]RuleHandler, len(config.Rules))
	for ruleName, rule := range config.Rules {

		ruleHandlers = append(ruleHandlers, RuleHandler{client, ruleName, rule, channels, nil, config.DryRun})
	}

	go func() {
		for t := range ticker.C {

			for i, handler := range ruleHandlers {
				handler.handle(t)
				ruleHandlers[i] = handler
			}
		}
	}()

	time.Sleep(maxDuration) // TODO find a better way to block here
	ticker.Stop()
	log.Println("Ticker stopped, exiting program")

	os.Exit(0)
}

// Log error message and exit program
func logFatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%v: %v", msg, err)
	}
}
