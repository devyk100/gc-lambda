package main

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"gc.yashk.dev/env"
	"gc.yashk.dev/gc_middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadImageFromUrlRequest struct {
	Token string `json:"token"`
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

	s3Client, err := gc_middleware.InitS3(ctx)
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
	return env.CLOUDFRONT_URL + filename.String(), nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	image_url := "https://farm3.staticflickr.com/2907/14066682001_88ce7553e4.jpg"
	// var greeting string
	// sourceIP := request.RequestContext.Identity.SourceIP

	// if sourceIP == "" {
	// 	greeting = "Hello, world!\n"
	// } else {
	// 	greeting = fmt.Sprintf("Hello, %s!\n", sourceIP)
	// }
	url, err := UploadFromUrlToS3(image_url)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       url,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
