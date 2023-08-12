package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func buildMessage(e CloudtrailEvent) slack.Message {

	// Header Section
	headerText := slack.NewTextBlockObject("plain_text", ":mag: Audit event detected :mag:", false, false)
	headerSection := slack.NewHeaderBlock(headerText, slack.HeaderBlockOptionBlockID("test_block"))

	reqParams := toJsonString(e.RequestParameters)
	log.WithField("data", reqParams).Debug("request_parameters")

	// Fields
	serviceField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Service:*\n%s", e.EventSource), false, false)
	whenField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Time:*\n%s", e.EventTime.Format("2006-01-02 15:04:05")), false, false)
	sourceIpField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Source IP:*\n%s", e.SourceIPAddress), false, false)
	actionField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Action:*\n%s", e.EventName), false, false)
	agentField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*UserAgent:*\n%s", e.UserAgent), false, false)
	regionField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Region:*\n%s", e.AwsRegion), false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, serviceField)
	fieldSlice = append(fieldSlice, whenField)
	fieldSlice = append(fieldSlice, sourceIpField)
	fieldSlice = append(fieldSlice, agentField)
	fieldSlice = append(fieldSlice, regionField)
	fieldSlice = append(fieldSlice, actionField)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	if e.UserIdentity.Arn == "" {
		e.UserIdentity.Arn = "N/A"
	}
	identityText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*User ARN:*\n```%s```", e.UserIdentity.Arn), false, false)
	identitySection := slack.NewSectionBlock(identityText, nil, nil)

	detailsText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Request Parameters:*\n```%v```", reqParams), false, false)
	detailsSection := slack.NewSectionBlock(detailsText, nil, nil)

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
		identitySection,
		detailsSection,
	)

	return msg
}

func toJsonString(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.WithError(err).Error("failed_to_parse_params")
	}
	return string(b)
}
