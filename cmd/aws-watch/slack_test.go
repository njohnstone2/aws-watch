package main

import (
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

var cloudtrailEvent = CloudtrailEvent{
	EventVersion: "1.08",
	UserIdentity: UserIdentity{
		Type:        "IAMUser",
		PrincipalID: "AIDA1AAAAAAAAA1AA1AAA",
		Arn:         "arn:aws:iam::123456789012:user/example",
		AccountID:   "123456789012",
		AccessKeyID: "ASIA1AAAAAAAAAAAA1AA",
		UserName:    "example",
	},
	EventTime:         time.Time{},
	EventSource:       "ec2.amazonaws.com",
	EventName:         "CreateSecurityGroup",
	AwsRegion:         "us-east-1",
	SourceIPAddress:   "1.2.3.4",
	UserAgent:         "AWS Internal",
	RequestParameters: "{\"groupName\":\"temp\",\"groupDescription\":\"example audit event\",\"vpcId\":\"vpc-a11a1111\"}",
	ResponseElements:  "{\"requestId\":\"add98acc-7c7d-44a0-829d-afef16381eff\",\"_return\":true,\"groupId\":\"sg-1234567890abcdefg\"}",
}

func TestBuildMessage(t *testing.T) {
	t.Run("Build Cloudtrail Message", func(t *testing.T) {

		msg := buildMessage(cloudtrailEvent)

		assert.Equal(t, len(msg.Blocks.BlockSet), 4)
		assert.Equal(t, msg.Blocks.BlockSet[0].BlockType(), slack.MBTHeader)
		assert.Equal(t, msg.Blocks.BlockSet[1].BlockType(), slack.MBTSection)
		assert.Equal(t, msg.Blocks.BlockSet[2].BlockType(), slack.MBTSection)
		assert.Equal(t, msg.Blocks.BlockSet[3].BlockType(), slack.MBTSection)
	})
}
