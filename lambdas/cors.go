package lambdas

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

var corsHeaders = map[string]string{
	"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "OPTIONS,POST,GET,PUT,DELETE",
}

func applyCORS(response events.APIGatewayProxyResponse) events.APIGatewayProxyResponse {
	if len(response.Headers) == 0 {
		response.Headers = make(map[string]string)
	}

	for k, v := range corsHeaders {
		response.Headers[k] = v
	}
	return response
}

type lambdaFunction func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func CORSWrapper(f lambdaFunction) lambdaFunction {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		response, err := f(ctx, request)
		response = applyCORS(response)
		return response, err
	}
}
