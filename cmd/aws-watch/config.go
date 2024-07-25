package main

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Subscribers []Subscriber `yaml:"subscribers"`
}

type Subscriber struct {
	Id        string    `yaml:"string"`
	Name      string    `yaml:"name"`
	Notifiers Notifiers `yaml:"notifiers"`
}

type Notifiers struct {
	GrafanaOncall GrafanaOncallConfig `yaml:"grafana-oncall"`
	Slack         SlackConfig         `yaml:"slack"`
}

type SlackConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Sources   []string `yaml:"sources"`
	ChannelId string   `yaml:"channel-id"`
}

type GrafanaOncallConfig struct {
	Enabled    bool     `yaml:"enabled"`
	WebhookUrl string   `yaml:"webhook-url"`
	Sources    []string `yaml:"sources"`
}

func NewConfig(path string) (*AppConfig, error) {

	var c = &AppConfig{}
	config, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(config, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func ReadConfig(bucket, filename, region string) (*AppConfig, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// sess, err := session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable, // Must be set to enable
	// 	Profile:           "default_temp",
	// })
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	rawObject, err := svc.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})

	if err != nil {
		return nil, err
	}

	var c = &AppConfig{}
	body, err := io.ReadAll(rawObject.Body)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(body, c)
	if err != nil {
		return c, err
	}

	return c, nil
}
