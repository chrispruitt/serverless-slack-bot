NAME=ssbot
FUNCTION_NAME=bot
VERSION=latest
DATE=`date +"%Y%m%d_%H%M%S"`
TEST_JSON='{"name": "My Poll"}'
STEP_FUNCTION_JSON='{"name": "My Step Execution Poll"}'

clean:
	rm -rf dist

updateLambda: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/main main.go
	cd dist && zip main.zip main
	aws lambda update-function-code --function-name ${FUNCTION_NAME} --zip-file fileb://${pwd}dist/main.zip

invoke:
	aws lambda invoke \
		--function-name "${FUNCTION_NAME}" \
		--log-type "Tail" \
		--payload $(TEST_JSON) \
		output/$(DATE).log \
		| jq -r '.LogResult' | base64 -D

invokeStepFunctionExecution:
	aws stepfunctions start-execution \
		--state-machine-arn arn:aws:states:us-east-1:770171891064:stateMachine:TriggerLambdaBot \
		--input $(STEP_FUNCTION_JSON)

