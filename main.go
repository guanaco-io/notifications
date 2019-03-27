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

	fmt.Printf("Starting notifications client\n")

	initError := validate()
	if initError != nil {
		log.Fatalf("Error initializing program: %v", initError)
	}

	ticker := time.NewTicker(Configuration.Alerta.ReloadInterval * time.Second)
	go func() {
		for t := range ticker.C {

			openAlerts := countOpenAlerts()
			log.Printf("Fetched %v open alerts from Alerta at %v\n", openAlerts, t)
		}
	}()

	time.Sleep(maxDuration) // TODO find a better way to block here
	ticker.Stop()
	log.Println("Ticker stopped")

	os.Exit(0)
}
