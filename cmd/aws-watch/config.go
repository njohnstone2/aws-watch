package main

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

type ConfigClient struct {
	client *s3.Client
}

func NewConfigClient(ctx context.Context, region string) (*ConfigClient, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &ConfigClient{
		client: client,
	}, nil
}

type S3GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func GetObjectFromS3(ctx context.Context, api S3GetObjectAPI, bucket, key string) ([]byte, error) {
	object, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()

	return io.ReadAll(object.Body)
}

func (c *ConfigClient) LoadConfig(ctx context.Context, bucket, filename string) (*AppConfig, error) {
	content, err := GetObjectFromS3(ctx, c.client, bucket, filename)
	if err != nil {
		return nil, err
	}

	var cfg = &AppConfig{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
