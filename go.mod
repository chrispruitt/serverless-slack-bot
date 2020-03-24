module github.com/chrispruitt/ssbot

go 1.13

replace github.com/chrispruitt/ssbot/bot => ./bot

require (
	github.com/aws/aws-lambda-go v1.16.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.6.0
	github.com/gin-gonic/gin v1.6.3
	github.com/slack-go/slack v0.6.4
)
