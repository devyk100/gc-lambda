package main

import (
	"context"
	"encoding/json"
	"log"

	db "gc.yashk.dev/checkemail/db_driver"
	"gc.yashk.dev/checkemail/env"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	dsn := env.DATABASE_URL

	// unmarshalling the hjson from the request
	var EmailCheckPayload EmailCheckPayload
	if err := json.Unmarshal([]byte(request.Body), &EmailCheckPayload); err != nil {
		log.Println("Failed to unmarshal the JSON")
	}

	// pgx pool to connect to the db
	var IsEmailTaken IsEmailTaken
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}
	defer pool.Close()

	// the query client to run the queries
	queries := db.New(pool)

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
