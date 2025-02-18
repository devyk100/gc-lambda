package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gc.yashk.dev/lambda/internal/env"
	"gc.yashk.dev/lambda/internal/middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadImageFromUrlRequest struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

type UploadImageFromUrlResponse struct {
	Url string `json:"url"`
}

func UploadFromUrlToS3(Url string) (string, error) {
	ctx := context.TODO()
	resp, err := http.Get(Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}

	s3Client, err := middleware.InitS3(ctx)
	if err != nil {
		return "", err
	}

	filename := uuid.New()
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(env.BucketName),
		Key:         aws.String(filename.String()),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(resp.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}
	return env.CLOUDFRONT_URL + filename.String(), nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var uploadImageFromUrlRequest UploadImageFromUrlRequest
	if err := json.Unmarshal([]byte(request.Body), &uploadImageFromUrlRequest); err != nil {
		log.Println("Failed to unmarshal the JSON")
	}

	var uploadImageFromUrlResponse UploadImageFromUrlResponse

	url, err := UploadFromUrlToS3(uploadImageFromUrlRequest.Url)
	uploadImageFromUrlResponse.Url = url
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	body, _ := json.Marshal(uploadImageFromUrlResponse)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
