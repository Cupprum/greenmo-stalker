package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Access query parameters
	lat1 := request.QueryStringParameters["lat1"]

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: `{"message": "Hello from Go! Lat1 is: ` + lat1 + `"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
