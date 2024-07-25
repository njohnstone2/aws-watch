package main

import (
	"context"
	"errors"

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

func (c *SecretsClient) GetAwsSecrets(secrets ...string) (map[string]string, error) {
	input := &secretsmanager.BatchGetSecretValueInput{
		SecretIdList: secrets,
	}

	result, err := c.client.BatchGetSecretValue(context.Background(), input)
	if err != nil {
		return map[string]string{}, err
	}

	if len(result.SecretValues) == 0 {
		return map[string]string{}, errors.New("no secrets found")
	}

	var values = make(map[string]string)
	for _, v := range result.SecretValues {
		values[*v.Name] = *v.SecretString
	}

	// Returns the the decrypted secrets as key/value pairs
	return values, nil
}
