package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//Handler is needed to be called by lambda
func handleRequest(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	responseBody := req.Headers["x-forwarded-for"] + "\n"
	response := events.ALBTargetGroupResponse{Body: responseBody, StatusCode: 200, IsBase64Encoded: false, StatusDescription: "200 OK", Headers: map[string]string{}}
	response.Headers["Content-Type"] = "text/html; charset=utf-8"
	prettyJSON, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		fmt.Printf("Could not unmarshal json. Error:\n%v", err)
	}
	fmt.Printf("Request:\n%s\n", string(prettyJSON))
	fmt.Printf("RemoteAddr: %v\n", responseBody)
	fmt.Printf("response: %v\n", response)
	return response, nil
}

func main() {
	fmt.Printf("Starting Handler\n")
	lambda.Start(handleRequest)
}
