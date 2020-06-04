# Slack Command Proxy

A proxy for Slack [slash commands](https://api.slack.com/interactivity/slash-commands) that does 
[request validation](https://api.slack.com/authentication/verifying-requests-from-slack) and publishes the payload to
[Cloud Pub/Sub](https://cloud.google.com/pubsub).

## Using
```shell script
git clone https://github.com/admiralobvious/slack-command-proxy
cd slack-command-proxy
```

1. Modify `config.json` with your own settings.

2. Create a Pub/Sub topic following this pattern:
```shell script
gcloud pubsub topics create slack-command-proxy-myteam-mycommand
```

3. Deploy Slack Command Proxy:
```shell script
gcloud functions deploy SlackCommandProxy --runtime go113 --trigger-http --set-env-vars "GCP_PROJECT=your-project-id" --allow-unauthenticated
```

## Why
- You work in an environment with more than one slash command
- You're already hosting your slash commands on GCP
- Your slash commands are implemented in more than one programming language
- Your slash commands are usually simple enough that you don't want to bundle a full-fledged Slack API library in your 
code just to do request validation
- Some of your slash commands may be written by less experienced (or who aren't primarily) developers and you'd rather 
they don't have to deal with request validation
- You don't wanna have to add a new service, open incoming ports and so on every time you add a new slash command
- You want the incoming requests to have a persistence layer

## Credits
[GCP Go Samples - functions](https://github.com/GoogleCloudPlatform/golang-samples/tree/9ca9b3f27ce69c46685ea66c70acc8a44815c56a/functions/slack)

[GCP Go Samples - pubsub](https://github.com/GoogleCloudPlatform/golang-samples/tree/9ca9b3f27ce69c46685ea66c70acc8a44815c56a/pubsub)
