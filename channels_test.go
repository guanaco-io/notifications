package main

import (
	"log"
	"testing"
)

func TestLoadChannelsConfig(t *testing.T) {

	Configuration, configError := Load("config/config.yaml")
	if configError != nil {
		t.Fatalf("cannot unmarshal data: %v", configError)
	}

	Channels, channelLoadError := LoadChannels(Configuration)
	if channelLoadError != nil {
		t.Fatalf("cannot load channels from config: %v", channelLoadError)
	}

	if len(Channels) != 3 {
		t.Fatalf("expected 3 channels")
	}
	for name, channel := range Channels {
		log.Printf("channel[%v]: %T %v", name, channel, channel)
	}
}
