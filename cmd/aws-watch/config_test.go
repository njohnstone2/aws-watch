package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

type MockConfigClient struct {
	GetObjectFunc  func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	LoadConfigFunc func(ctx context.Context, bucket, filename string) (*AppConfig, error)
}

func (m *MockConfigClient) LoadConfig(ctx context.Context, bucket, filename string) (*AppConfig, error) {
	return m.LoadConfigFunc(ctx, bucket, filename)
}

func TestLoadConfig(t *testing.T) {

	cases := []struct {
		cfgClient MockConfigClient
		bucket    string
		key       string
		expect    *AppConfig
	}{
		{
			cfgClient: MockConfigClient{
				LoadConfigFunc: func(ctx context.Context, bucket string, filename string) (*AppConfig, error) {
					return &AppConfig{
						Subscribers: []Subscriber{
							{
								Id:   "teamA",
								Name: "Team A",
								Notifiers: Notifiers{
									GrafanaOncall: GrafanaOncallConfig{
										Enabled:    true,
										WebhookUrl: "http://localhost:8080",
										Sources:    []string{"ec2"},
									},
									Slack: SlackConfig{
										Enabled:   true,
										Sources:   []string{"iam"},
										ChannelId: "AAAAAAAAAAA",
									},
								},
							},
							{
								Id:   "teamB",
								Name: "Team B",
								Notifiers: Notifiers{
									Slack: SlackConfig{
										Enabled:   true,
										Sources:   []string{"*"},
										ChannelId: "BBBBBBBBBBB",
									},
								},
							},
						},
					}, nil
				},
			},
			bucket: "fooBucket",
			key:    "barKey",
			expect: &AppConfig{
				Subscribers: []Subscriber{
					{
						Id:   "teamA",
						Name: "Team A",
						Notifiers: Notifiers{
							GrafanaOncall: GrafanaOncallConfig{
								Enabled:    true,
								WebhookUrl: "http://localhost:8080",
								Sources:    []string{"ec2"},
							},
							Slack: SlackConfig{
								Enabled:   true,
								Sources:   []string{"iam"},
								ChannelId: "AAAAAAAAAAA",
							},
						},
					},
					{
						Id:   "teamB",
						Name: "Team B",
						Notifiers: Notifiers{
							Slack: SlackConfig{
								Enabled:   true,
								Sources:   []string{"*"},
								ChannelId: "BBBBBBBBBBB",
							},
						},
					},
				},
			},
		},
		{
			cfgClient: MockConfigClient{
				LoadConfigFunc: func(ctx context.Context, bucket string, filename string) (*AppConfig, error) {
					return &AppConfig{
						Subscribers: []Subscriber{},
					}, nil
				},
			},
			bucket: "secondBucket",
			key:    "secondKey",
			expect: &AppConfig{
				Subscribers: []Subscriber{},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.Background()

			actual, err := tt.cfgClient.LoadConfig(ctx, tt.bucket, tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
