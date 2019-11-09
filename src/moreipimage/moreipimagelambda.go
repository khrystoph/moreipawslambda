//lambda image handler/server
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	awsRegion = "us-west-2"
)

var (
	s3bucketname = os.Getenv("BUCKET")
	bucketKey    = os.Getenv("IMAGE")
	sess         *session.Session
	awsConfig    = &aws.Config{
		Region: aws.String(awsRegion),
	}
	getObjectInput = &s3.GetObjectInput{
		Bucket: aws.String(s3bucketname),
		Key:    aws.String(bucketKey),
	}
	listObjectsInput = &s3.ListObjectsV2Input{
		Bucket: aws.String(s3bucketname),
	}
)

func init() {
	sess = session.Must(session.NewSession(awsConfig))
}

//handleRequest is needed to be called by lambda
func handleRequest(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	responseBody, err := createBase64Image()
	if err != nil {
		fmt.Printf("Unable to create base64String:\n%s", err)
		return events.ALBTargetGroupResponse{
			Body:              "Error pulling image",
			StatusCode:        500,
			IsBase64Encoded:   false,
			StatusDescription: "500 ERROR",
			Headers:           map[string]string{},
		}, err
	}
	response := events.ALBTargetGroupResponse{
		Body:              responseBody,
		StatusCode:        200,
		IsBase64Encoded:   true,
		StatusDescription: "200 OK",
		Headers:           map[string]string{},
	}
	response.Headers["Content-Type"] = "image/jpeg;"
	fmt.Printf("Request: %v\n", req)
	fmt.Printf("RemoteAddr: %v\n", responseBody)
	fmt.Printf("response: %v\n", response)
	return response, nil
}

func createBase64Image() (image string, err error) {
	s3ImageObject, err := pullImage()
	if err != nil {
		fmt.Printf("Unable to pull image. Error message:\n%s", err)
		return "", err
	}
	fmt.Println("converting to base64")
	image, err = convertBase64(s3ImageObject)
	fmt.Println("returning image to main function")
	return image, err
}

func convertBase64(s3Image *s3.GetObjectOutput) (encodedString string, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(s3Image.Body)
	encodedString = base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Printf("%v\n", encodedString)
	return encodedString, nil
}

func pullImage() (*s3.GetObjectOutput, error) {
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Println("Error creating session.")
		return nil, err
	}

	fmt.Printf("getObjectInput:\n%v\n", *getObjectInput)
	result, err := s3.New(sess).GetObject(getObjectInput)
	fmt.Println("checking for errors from GetObject.")
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return result, err
	}
	fmt.Printf("returning result from GetObject\n")
	return result, nil
}

func main() {
	fmt.Printf("Starting Handler\n")
	lambda.Start(handleRequest)
}
