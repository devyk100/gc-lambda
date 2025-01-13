package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gc.yashk.dev/env"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/s3/actions"
	"github.com/golang-jwt/jwt/v5"
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

func Initialize(fileName string) (string, error) {
	// unsure of the context here
	ctx := context.TODO()

	// loading from the default config, other direct structs do not work here
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(env.Key, env.Secret, ""),
		))
	if err != nil {
		return "", err
	}
	// the s3 client from function
	s3Client := s3.NewFromConfig(cfg)
	// for presigning purposes we need a seperate presigning client
	presignClient := s3.NewPresignClient(s3Client)
	presigner := actions.Presigner{PresignClient: presignClient}
	// details of the operation, the url and methods allowed
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

	// parsing and checking the user from this jwt
	fmt.Println("The token is", putImageUrlRequest.Token)
	t, err := jwt.Parse(putImageUrlRequest.Token, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(env.NEXTAUTH_SECRET)
		return secret, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	if !t.Valid {
		fmt.Println("THE TOKEN WAS NOT VALIDATED MAN")
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 403,
		}, nil
	}

	// make the database queries, to save this image's keys and all associate it with this user / email
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Could not parse claims")
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 403,
		}, nil
	}
	email, emailExists := claims["email"].(string)
	if !emailExists {
		fmt.Println("Email claim not found")
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 403,
		}, nil
	}

	fmt.Println(email, "is the email") // do SOME DB CALLS NOW

	fileName := uuid.New()
	putUrl, err := Initialize(fileName.String())
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
