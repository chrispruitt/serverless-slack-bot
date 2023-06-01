module github.com/chrispruitt/serverless-slack-bot

go 1.16

replace github.com/chrispruitt/serverless-slack-bot/bot => ./bot

require (
	github.com/aws/aws-lambda-go v1.26.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.11.0
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/gin-gonic/gin v1.9.1
	github.com/mattes/go-asciibot v0.0.0-20190603170252-3fa6d766c482
	github.com/slack-go/slack v0.9.4
)
