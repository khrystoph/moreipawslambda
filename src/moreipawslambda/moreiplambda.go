package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//Handler is needed to be called by lambda
func handleRequest(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	responseBody := req.Headers["x-forwarded-for"] + "\n"
	response := events.ALBTargetGroupResponse{Body: responseBody, StatusCode: 200, IsBase64Encoded: false, StatusDescription: "200 OK", Headers: map[string]string{}}
	response.Headers["Content-Type"] = "text/html; charset=utf-8"
	fmt.Printf("Request: %v\n", req)
	fmt.Printf("RemoteAddr: %v\n", responseBody)
	fmt.Printf("response: %v\n", response)
	return response, nil
}

func main() {
	fmt.Printf("Starting Handler\n")
	lambda.Start(handleRequest)
}
