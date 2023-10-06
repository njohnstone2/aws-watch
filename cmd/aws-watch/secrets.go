package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsClient struct {
	client *secretsmanager.Client
}

func NewSecretsClient(region string) *SecretsClient {
	config, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	return &SecretsClient{
		client: svc,
	}
}

func (c *SecretsClient) GetAwsSecret(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := c.client.GetSecretValue(context.Background(), input)
	if err != nil {
		return "", err
	}

	// Returns the decrypted secret
	return *result.SecretString, nil
}
