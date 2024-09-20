package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const (
	secret_slack_token              = "slack_token"
	cloudwatch_eks_audit_log_prefix = "kube-apiserver-audit-"
)

func handler(request events.CloudwatchLogsEvent) error {
	LOG_LEVEL := os.Getenv("LOG_LEVEL")
	AWS_REGION := os.Getenv("REGION")
	S3_BUCKET_NAME := os.Getenv("S3_BUCKET_NAME")
	S3_FILENAME := os.Getenv("S3_FILENAME")
	setLogger(LOG_LEVEL)

	log.WithFields(log.Fields{
		"data":       request.AWSLogs.Data,
		"log_level":  LOG_LEVEL,
		"aws_region": AWS_REGION,
	}).Info("inputs")

	ctx := context.Background()

	// Load Subscriber Configuration
	configClient, err := NewConfigClient(ctx, AWS_REGION)
	if err != nil {
		log.WithError(err).Error("failed_to_init_s3_client")
		return err
	}

	c, cErr := configClient.LoadConfig(ctx, S3_BUCKET_NAME, S3_FILENAME)
	if cErr != nil {
		log.WithError(cErr).Error("failed_to_load_config")
		return err
	}

	if len(c.Subscribers) == 0 {
		err := fmt.Errorf("one or more subscribers must be configured")
		log.WithError(err).Fatal("No subscribers found. Check config file exists.")
		return err
	}

	// fetch secrets
	secretsClient := NewSecretsClient(AWS_REGION)

	slackToken, err := secretsClient.GetAwsSecret(secret_slack_token)
	if err != nil {
		log.WithError(err).Error("failed_to_get_secret")
		return err
	}

	oncallClient := NewGrafanaOncallClient()

	parsed, err := request.AWSLogs.Parse()
	if err != nil {
		log.WithError(err).Error("failed_to_parse_event")
		return err
	}

	log.WithField("count", len(parsed.LogEvents)).Info("events_received")
	if len(parsed.LogEvents) > 0 {
		for _, v := range parsed.LogEvents {
			log.WithField("data", v.Message).Debug("parsed_message")

			var event Event
			if strings.HasPrefix(parsed.LogStream, cloudwatch_eks_audit_log_prefix) {
				eksEvent, eksErr := parseEKSEvent(v)
				if eksErr != nil {
					return eksErr
				}

				event.EventSource = "EKS"
				event.EKS = eksEvent
				event.EKS.ClusterName = getEKSClusterName(parsed.LogGroup)

				log.WithFields(log.Fields{
					"source": event.EventSource,
					"uri":    event.EKS.RequestURI,
				}).Info("eks_event")
			} else {
				ctEvent, ctErr := parseCloudtrailEvent(v)
				if ctErr != nil {
					return ctErr
				}

				event.EventSource = ctEvent.EventSource
				event.Cloudtrail = ctEvent

				log.WithFields(log.Fields{
					"source": event.EventSource,
					"arn":    event.Cloudtrail.UserIdentity.Arn,
				}).Info("cloudtrail_event")
			}

			// Iterate over configured subscribers
			for _, s := range c.Subscribers {
				log.WithFields(log.Fields{
					"team":           s.Name,
					"event_source":   event.EventSource,
					"oncall_sources": s.Notifiers.GrafanaOncall.Sources,
				}).Info("Evaluating Subscriber...")
				if sliceContains(s.Notifiers.GrafanaOncall.Sources, event.EventSource) {
					msg := oncallClient.buildMessage(event)
					uid := uuid.New()
					alert := &Alert{
						UID:     uid.String(),
						Title:   fmt.Sprintf("Audit Event (%s)", event.EventSource),
						State:   "alerting",
						Message: msg,
					}
					log.WithFields(log.Fields{
						"team":        s.Name,
						"source":      event.EventSource,
						"webhook_url": s.Notifiers.GrafanaOncall.WebhookUrl,
						"msg":         msg,
					}).Info("Posting event to Grafana Oncall")

					err := oncallClient.CreateAlert(s.Notifiers.GrafanaOncall.WebhookUrl, alert)
					if err != nil {
						log.WithError(err).WithFields(log.Fields{
							"team":        s.Name,
							"source":      event.EventSource,
							"webhook_url": s.Notifiers.GrafanaOncall.WebhookUrl,
						}).Error("failed_post_to_grafana_oncall")
					}
					log.WithFields(log.Fields{
						"webhook_url": s.Notifiers.GrafanaOncall.WebhookUrl,
						"source":      event.EventSource,
						"alert_title": alert.Title,
					}).Info("Message successfully sent to Grafana Oncall")
				}

				if sliceContains(s.Notifiers.Slack.Sources, event.EventSource) {
					log.WithFields(log.Fields{
						"team":   s.Name,
						"source": event.EventSource,
					}).Debug("Posting event to Slack")

					msg := buildMessage(event)
					pErr := slackPost(slackToken, s.Notifiers.Slack.ChannelId, msg)
					if pErr != nil {
						log.WithFields(log.Fields{
							"channel_id": s.Notifiers.Slack.ChannelId,
						}).WithError(pErr).Error("failed_post_to_slack")
					}
				}
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

func parseEKSEvent(e events.CloudwatchLogsLogEvent) (*EKSEvent, error) {
	var event EKSEvent
	err := json.Unmarshal([]byte(e.Message), &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func parseCloudtrailEvent(e events.CloudwatchLogsLogEvent) (*CloudtrailEvent, error) {
	var event CloudtrailEvent
	err := json.Unmarshal([]byte(e.Message), &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func setLogger(level string) {
	log.SetFormatter(&log.JSONFormatter{})

	switch strings.ToUpper(level) {
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

// checks for a prefix match in a slice of strings or is a wildcard match
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if strings.HasPrefix(e, a) || a == "*" {
			return true
		}
	}
	return false
}

func getEKSClusterName(s string) string {
	clusterName := s
	clusterName = strings.TrimPrefix(clusterName, "/aws/eks/")
	clusterName = strings.TrimSuffix(clusterName, "/cluster")

	return clusterName
}
