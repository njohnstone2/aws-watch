# Slack Bot Setup

## Prerequisites
Ensure you have access to the slack workspace that will host your slack bot

## Steps
1. Navigate to https://api.slack.com/apps
2. Select `Create New App` > `From Scratch`
3. Set your app name and desired workspace
4. Select `Bots` > `Review Scopes to Add`
5. Under `Scopes` add the following to `Bot Token Scopes`:
- chat:write
6. Under `OAuth Tokens for Your Workspace` select `Install to Workspace`
7. Take note of the generated bot token as you will need to add it to your lambda configuration \
NOTE: the bot token should begin with `xoxb-`

## Adding the bot to a channel / Finding your channel ID
1. Navigate to https://app.slack.com/
2. Select the slack channel and mention the bot to add it (e.g. @AWSWatch)
3. Take note of the channel ID as you will need to add it to your lambda configuration \
NOTE: The channel ID is the last section of the url (e.g. `https://app.slack.com/client/<WORKSPACE_ID>/<CHANNEL_ID>`)
