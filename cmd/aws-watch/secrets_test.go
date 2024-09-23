package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSecretsClient struct {
	GetAwsSecretFunc  func(secretName string) (string, error)
	GetAwsSecretsFunc func(secrets ...string) (map[string]string, error)
}

func (m MockSecretsClient) GetAwsSecret(secretName string) (string, error) {
	return m.GetAwsSecretFunc(secretName)
}

func (m MockSecretsClient) GetAwsSecrets(secrets ...string) (map[string]string, error) {
	return m.GetAwsSecretsFunc(secrets...)
}

const AWS_REGION = "us-east-1"

func TestGetSecret(t *testing.T) {
	t.Run("Retrieve a single secret from AWS SecretsManager", func(t *testing.T) {

		secretsClient := &MockSecretsClient{
			GetAwsSecretFunc: func(secrets string) (string, error) {
				return "secretValue", nil
			},
		}

		expected := "secretValue"
		actual, err := secretsClient.GetAwsSecret("slack_token")

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestBatchGetSecrets(t *testing.T) {
	t.Run("Retrieve multiple secrets from AWS SecretsManager", func(t *testing.T) {

		secretsClient := &MockSecretsClient{
			GetAwsSecretsFunc: func(secrets ...string) (map[string]string, error) {
				return map[string]string{
					"slack_token":            "abc",
					"slack_channel_id_teamA": "def",
					"slack_channel_id_teamB": "123",
				}, nil
			},
		}

		// secrets, err := secretsClient.GetAwsSecrets("slack_token", "slack_channel_id_teamA", "slack_channel_id_teamB")
		secrets, err := secretsClient.GetAwsSecrets()
		assert.NoError(t, err)
		assert.Equal(t, 3, len(secrets))
	})
}

func TestSecretsError(t *testing.T) {
	t.Run("Fetch Invalid Secret Name", func(t *testing.T) {

		secretsClient := NewSecretsClient(AWS_REGION)
		_, err := secretsClient.GetAwsSecret("invalid")

		assert.Error(t, err)
	})
}
