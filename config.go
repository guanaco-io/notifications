package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Config struct {
	DryRun bool `yaml:"dry_run"`

	Alerta          Alerta                   `yaml:"alerta"`
	ChannelSettings ChannelSettings          `yaml:"channel_settings"`
	Channels        map[string]ChannelConfig `yaml:"channels"`
	Rules           map[string]Rule          `yaml:"rules"`
}

type Alerta struct {
	Endpoint       string        `yaml:"endpoint"`
	Webui          string        `yaml:"webui"`
	ReloadInterval time.Duration `yaml:"reload_interval"`
}

type ChannelSettings struct {
	Slack Slack `yaml:"slack"`
	Smtp  Smtp  `yaml:"smtp"`
}

type Slack struct {
	WebhookUrl string `yaml:"webhook_url"`
}

type Smtp struct {
	Server    string `yaml:"server"`
	Port      int    `yaml:"port"` // 465 for ssl, 587 for non-ssl
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	From      string `yaml:"from"`
	FromName  string `yaml:"from_name"`
	Anonymous bool   `yaml:"anonymous"`
	Ssl       bool   `yaml:"ssl"`
}

type ChannelConfig struct {
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

type Rule struct {
	Filter   string   `yaml:"filter"`
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

	return config, nil
}
