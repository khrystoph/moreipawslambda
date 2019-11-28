package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var (
	inputRequest = events.ALBTargetGroupRequest{
		HTTPMethod: "GET",
		Headers: map[string]string{
			"x-forwarded-for":   "1.2.3.4",
			"x-forwarded-port":  "80",
			"x-forwarded-proto": "http",
			"host":              "2.2.2.2",
		},
	}
	requestContext context.Context = nil
)

func TestHandleRequest(t *testing.T) {
	response, err := handleRequest(requestContext, inputRequest)
	assert.Nil(t, err, "There should not be an error returned.")
	assert.Equal(t, "1.2.3.4\n", response.Body, "Expecting 1.2.3.4\\n for returned IP.")
}
