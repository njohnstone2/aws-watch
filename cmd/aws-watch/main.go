package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const (
	secret_slack_token      = "slack_token"
	secret_slack_channel_id = "slack_channel_id"
)

func handler(request events.CloudwatchLogsEvent) error {
	setLogger()
	LOG_LEVEL := os.Getenv("LOG_LEVEL")
	AWS_REGION := os.Getenv("REGION")

	log.WithFields(log.Fields{
		"data":       request.AWSLogs.Data,
		"log_level":  LOG_LEVEL,
		"aws_region": AWS_REGION,
	}).Debug("inputs")

	// fetch secrets
	secretsClient := NewSecretsClient(AWS_REGION)

	slackToken, err := secretsClient.GetAwsSecret(secret_slack_token)
	if err != nil {
		log.WithError(err).Error("failed_to_get_secret")
	}

	slackChannelId, channelErr := secretsClient.GetAwsSecret(secret_slack_channel_id)
	if err != nil {
		log.WithError(channelErr).Error("failed_to_get_secret")
	}

	log.WithFields(log.Fields{
		"token":      slackToken,
		"channel_id": slackChannelId,
	}).Debug("secrets")

	parsed, err := request.AWSLogs.Parse()
	if err != nil {
		log.WithError(err).Error("failed_to_parse_event")
		return err
	}

	log.WithField("count", len(parsed.LogEvents)).Info("events_received")
	if len(parsed.LogEvents) > 0 {
		for _, v := range parsed.LogEvents {
			log.WithField("data", v.Message).Debug("parsed_message")

			var event CloudtrailEvent
			err := json.Unmarshal([]byte(v.Message), &event)
			if err != nil {
				log.WithError(err).Error("failed_to_unmarshal_event")
				return err
			}

			msg := buildMessage(event)
			pErr := slackPost(slackToken, slackChannelId, msg)
			if pErr != nil {
				log.WithError(pErr).Error("failed_post_to_slack")
				return pErr
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

func setLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	loglevel := os.Getenv("LOG_LEVEL")

	switch strings.ToUpper(loglevel) {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func printMessagePayload(msg slack.Message) {
	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		log.WithError(err).Error("failed_marshalling")
	}
	log.WithField("data", string(b)).Info("message_payload")
}

func slackPost(token, channelId string, msg slack.Message) error {
	api := slack.New(token)
	channelID, timestamp, err := api.PostMessage(
		channelId,
		slack.MsgOptionBlocks(msg.Blocks.BlockSet...),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"channel_id": channelID,
		"timestamp":  timestamp,
	}).Info("Message successfully sent to channel")

	return nil
}
