package main

import (
	"context"
	"fmt"
	"time"

	"gc.yashk.dev/env"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/s3/actions"
)

type PutImageUrlRequest struct {
	Extension int    `json:"extension"`
	Token     string `json:"token"`
}

type PutImageUrlResponse struct {
	Url string `json:"url"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(env.Key, env.Secret, ""),
		))
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	s3Client := s3.NewFromConfig(cfg)
	// bucketBasics := actions.BucketBasics{S3Client: s3Client}
	presignClient := s3.NewPresignClient(s3Client)
	presigner := actions.Presigner{PresignClient: presignClient}
	hey := "fwmaofew"
	request1, err := presigner.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &env.BucketName,
		Key:    &hey,
	}, func(po *s3.PresignOptions) {
		po.Expires = time.Duration(40 * int64(time.Second))
	})
	if err != nil {
		fmt.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			"bucketName", " objectKey", err)
	}

	fmt.Println(request1.URL, request1)

	// sourceIP := request.RequestContext.Identity.SourceIP
	// if sourceIP == "" {
	// 	greeting := "Hello, world!\n"
	// } else {
	// 	greeting = fmt.Sprintf("Hello, %s!\n", sourceIP)
	// }

	return events.APIGatewayProxyResponse{
		Body:       request1.URL,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
