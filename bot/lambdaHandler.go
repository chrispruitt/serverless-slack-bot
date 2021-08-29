package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack/slackevents"
)

var ginLambda *ginadapter.GinLambdaV2
var r *gin.Engine

func StatusRoute(c *gin.Context) {
	c.JSON(200, "ok")
}

func SlackEventRoute(c *gin.Context) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	body := buf.String()

	fmt.Println("Attempting to verify")
	fmt.Println(json.RawMessage(body))
	fmt.Println("Attempting to verify")

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

func stripApiGatewayStageName(req events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPRequest {
	request := req
	request.RawPath = strings.Replace(request.RawPath, fmt.Sprintf("/%s", request.RequestContext.Stage), "", 1)
	request.RequestContext.HTTP.Path = strings.Replace(request.RequestContext.HTTP.Path, fmt.Sprintf("/%s", request.RequestContext.Stage), "", 1)
	return request
}

func apiGatewayProxyRequestHandler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	_, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	// Stip stage if exists from path
	req = stripApiGatewayStageName(req)

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

	var request events.APIGatewayV2HTTPRequest
	json.Unmarshal(payload, &request)

	// Check if event was actually an API Gateway Event by checking if RawPath exists
	if request.RawPath != "" {
		fmt.Println("Received an API Gateway Proxy Event.")
		response, err := apiGatewayProxyRequestHandler(ctx, request)
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
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r = gin.Default()

	r.POST("/slack-event", SlackEventRoute)
	r.GET("/status", StatusRoute)

	ginLambda = ginadapter.NewV2(r)

	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		handle := new(CustomHandle)
		lambda.StartHandler(handle)
	} else {
		r.Run()
	}
}
