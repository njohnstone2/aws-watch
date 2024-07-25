package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ExampleConfigPath = "../../config/example.yaml"
)

func TestValidateConfig(t *testing.T) {
	t.Run("Validate Grafana Oncall Config", func(t *testing.T) {

		c, err := NewConfig(ExampleConfigPath)

		assert.NoError(t, err)
		assert.Equal(t, len(c.Subscribers), 2)
		assert.Equal(t, c.Subscribers[0].Notifiers.GrafanaOncall.Enabled, true)
		assert.Equal(t, c.Subscribers[0].Notifiers.GrafanaOncall.WebhookUrl, "http://localhost:8000/asdf")
		assert.Equal(t, c.Subscribers[0].Notifiers.Slack.Enabled, false)
	})
}

func TestMissingConfig(t *testing.T) {
	t.Run("Validate Missing Config", func(t *testing.T) {

		_, err := NewConfig("")

		assert.Error(t, err)
	})
}
