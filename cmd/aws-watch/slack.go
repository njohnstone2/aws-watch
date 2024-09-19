package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func buildMessage(e Event) slack.Message {
	var msg = slack.Message{}
	switch e.EventSource {
	case "EKS":
		msg = buildEKSMessage(*e.EKS)
	default:
		msg = buildCloudTrailMessage(*e.Cloudtrail)
	}

	return msg
}

func buildEKSMessage(e EKSEvent) slack.Message {

	// Header Section
	headerText := slack.NewTextBlockObject("plain_text", ":mag: Audit event detected :mag:", false, false)
	headerSection := slack.NewHeaderBlock(headerText, slack.HeaderBlockOptionBlockID("test_block"))

	if e.User.Username == "" {
		e.User.Username = "N/A"
	}

	// Fields
	serviceField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Service:*\n%s", "EKS"), false, false)
	whenField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Time:*\n%s", e.RequestReceivedTimestamp.Format("2006-01-02 15:04:05")), false, false)
	sourceIpField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Source IP:*\n%s", e.SourceIPs[0]), false, false)
	actionField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Action:*\n%s", e.Verb), false, false)
	agentField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*UserAgent:*\n%s", e.UserAgent), false, false)
	userField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Username:*\n%s", e.User.Username), false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, serviceField)
	fieldSlice = append(fieldSlice, whenField)
	fieldSlice = append(fieldSlice, sourceIpField)
	fieldSlice = append(fieldSlice, agentField)
	fieldSlice = append(fieldSlice, actionField)
	fieldSlice = append(fieldSlice, userField)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	uriText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Request URI:*\n```%s```", e.RequestURI), false, false)
	uriSection := slack.NewSectionBlock(uriText, nil, nil)

	reqObj := toJsonString(e.RequestObject.Data)
	log.WithField("data", reqObj).Debug("request_object")

	reqText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Request Object:*\n```%v```", reqObj), false, false)
	requestSection := slack.NewSectionBlock(reqText, nil, nil)

	responseObject := toJsonString(e.ResponseObject.Data)
	log.WithField("data", responseObject).Debug("response_object")

	responseText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Response Object:*\n```%v```", responseObject), false, false)
	responseSection := slack.NewSectionBlock(responseText, nil, nil)

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
		uriSection,
		requestSection,
		responseSection,
	)

	return msg
}

func buildCloudTrailMessage(e CloudtrailEvent) slack.Message {

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
