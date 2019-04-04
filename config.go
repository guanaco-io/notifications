package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	DryRun bool

	Alerta Alerta `yaml:"alerta"`
	ChannelSettings ChannelSettings `yaml:"channel_settings"`
	Channels map[string]Channel `yaml:"channels"`
	Rules map[string]Rule `yaml:"rules"`
}

type Alerta struct {

	Endpoint string `yaml:"endpoint"`
	Webui string `yaml:"webui"`
	ReloadInterval time.Duration `yaml:"reload_interval"`
}

type ChannelSettings struct {

	Slack Slack `yaml:"slack"`
	Smtp Smtp `yaml:"smtp"`
}

type Slack struct {

	WebhookUrl string `yaml:"webhook_url"`
}

type Smtp struct {

	Server string `yaml:"server"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	From string `yaml:"from"`
	Ssl bool `yaml:"ssl"`
}

type Channel struct {

	Type string `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

type Rule struct {

	Filter string `yaml:"filter"`
	Channels []string `yaml:"channels"`
}

func Load(filename string) (Config, error) {

	var config Config

	data, readFileError := ioutil.ReadFile(filename)
	if readFileError != nil {
		return config, readFileError
	}

	unmarshallError := yaml.Unmarshal(data, &config)
	if unmarshallError != nil {
		return config, unmarshallError
	}

	config.DryRun = contains(os.Args, "--dry-run")

	return config, nil
}

func contains(slice []string, lookup string) bool {
	for _, element := range slice {
		if element == lookup {
			return true
		}
	}
	return false
}