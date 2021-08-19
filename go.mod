module github.com/chrispruitt/serverless-slack-bot

go 1.16

replace github.com/chrispruitt/serverless-slack-bot/bot => ./bot

require (
	github.com/aws/aws-lambda-go v1.26.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.11.0
	github.com/gin-gonic/gin v1.7.0
	github.com/slack-go/slack v0.9.4
)
