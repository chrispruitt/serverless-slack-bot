package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack/slackevents"
)

var ginLambda *ginadapter.GinLambda
var r *gin.Engine

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r = gin.Default()

	r.POST("/slack-event", SlackEventHandler)

	ginLambda = ginadapter.New(r)
}

func SlackEventHandler(c *gin.Context) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	body := buf.String()

	// Verify the request came from slack
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN")}))
	if e != nil {
		fmt.Println(e.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
		return
	}

	// Verify event URL when setting up bot
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, r.Challenge)
		return
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			HandleMentionEvent(ev)
		}
	}

	c.JSON(200, "OK")
}

func apiGatewayProxyRequestHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func PollEventHandler(ctx context.Context, req PollEvent) (PollEvent, error) {
	_, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	return PollEvent{
		Name:   req.Name,
		Status: "Success",
	}, nil
}

type CustomHandle struct{}

func (handler CustomHandle) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	fmt.Println("RAW REQUEST")
	fmt.Println(string(payload[:]))

	var apiGatewayProxyRequest events.APIGatewayProxyRequest
	json.Unmarshal(payload, &apiGatewayProxyRequest)

	// Check if event was actually an API Gateway Event by checking if HTTPMethod exists
	if apiGatewayProxyRequest.HTTPMethod != "" {
		fmt.Println("Received an API Gateway Proxy Event.")
		response, err := apiGatewayProxyRequestHandler(ctx, apiGatewayProxyRequest)
		if err != nil {
			return nil, err
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return responseBytes, nil
	}

	var pollEvent PollEvent
	json.Unmarshal(payload, &pollEvent)

	if pollEvent.Name != "" {
		fmt.Println("Received a Poll Request Event.")

		response, err := PollEventHandler(ctx, pollEvent)
		if err != nil {
			return nil, err
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return responseBytes, nil
	}

	return []byte(`{"error":"Unable to handle this event type."}`), nil
}

func Start() {
	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		handle := new(CustomHandle)
		lambda.StartHandler(handle)
	} else {
		r.Run()
	}
}
