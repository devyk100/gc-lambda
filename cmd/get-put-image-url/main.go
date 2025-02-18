package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gc.yashk.dev/lambda/internal/env"
	"gc.yashk.dev/lambda/internal/middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type PutImageUrlRequest struct {
	Extension string `json:"extension"`
	Token     string `json:"token"`
}

type PutImageUrlResponse struct {
	Url      string `json:"url"`
	FileName string `json:"filename"`
}

func GetSignedPutUrl(fileName string) (string, error) {
	// unsure of the context here
	ctx := context.TODO()
	presigner, err := middleware.InitS3Presigner(ctx)
	if err != nil {
		return "", err
	}

	imageUploadDetails, err := presigner.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &env.BucketName,
		Key:    &fileName,
	}, func(po *s3.PresignOptions) {
		po.Expires = time.Duration(40 * int64(time.Second))
	})
	if err != nil {
		return "", err
	}

	return imageUploadDetails.URL, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// sourceIP := request.RequestContext.Identity.SourceIP
	var putImageUrlRequest PutImageUrlRequest
	if err := json.Unmarshal([]byte(request.Body), &putImageUrlRequest); err != nil {
		fmt.Println("Failed to unmarshal", err.Error())
	}
	ctx := context.TODO()
	queries, pool, err := gc_middleware.InitDb(ctx)
	defer pool.Close()
	if err != nil {

	}
	isAuthenticated, err := gc_middleware.JwtAuth(ctx, &putImageUrlRequest.Token, queries)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       string("you are not authenticated"),
		}, nil
	}

	if !isAuthenticated {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       string("You are not authenticated"),
		}, nil
	}

	fileName := uuid.New()
	putUrl, err := GetSignedPutUrl(fileName.String())
	if err != nil {
		fmt.Println(err.Error())
	}
	var putImageUrlResponse PutImageUrlResponse
	putImageUrlResponse.Url = putUrl
	putImageUrlResponse.FileName = fileName.String()
	responseBody, err := json.Marshal(putImageUrlResponse)
	if err != nil {
		fmt.Println("error", err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
