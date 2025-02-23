package main

import (
	"context"
	"encoding/json"
	"log"

	"gc.yashk.dev/lambda/internal/middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
)

type IsEmailTaken struct {
	Success bool `json:"success"`
}

type EmailCheckPayload struct {
	Email string `json:"email"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// creating contexts and loading the db url
	ctx := context.Background()

	// unmarshalling the hjson from the request
	var EmailCheckPayload EmailCheckPayload
	if err := json.Unmarshal([]byte(request.Body), &EmailCheckPayload); err != nil {
		log.Println("Failed to unmarshal the JSON")
	}

	var IsEmailTaken IsEmailTaken

	queries, pool, err := middleware.InitDb(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}
	defer pool.Close()

	_, err = queries.GetUserFromEmail(ctx, EmailCheckPayload.Email)
	//error handling for this case
	if err != nil {
		// no rows found in the sql query
		if err == pgx.ErrNoRows {
			IsEmailTaken.Success = false
			body, _ := json.Marshal(IsEmailTaken)
			return events.APIGatewayProxyResponse{
				Body:       string(body),
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type":                 "application/json",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type, Authorization",
				},
			}, nil
		} else {
			// some other sql error
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 500,
				Headers: map[string]string{
					"Content-Type":                 "application/json",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type, Authorization",
				},
			}, nil
		}
	} else {
		// if a user was successfully found in this case
		IsEmailTaken.Success = true
		body, _ := json.Marshal(IsEmailTaken)
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
