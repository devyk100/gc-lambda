package main

import (
	"context"
	"encoding/json"
	"log"

	db "gc.yashk.dev/checkusername/db_driver"
	"gc.yashk.dev/checkusername/env"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IsUsernameAvailable struct {
	Success bool `json:"success"`
}

type UsernameCheckPayload struct {
	Username string `json:"username"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// creating contexts and loading the db url
	ctx := context.Background()
	dsn := env.DATABASE_URL

	// unmarshalling the hjson from the request
	var usernameCheckPayload UsernameCheckPayload
	if err := json.Unmarshal([]byte(request.Body), &usernameCheckPayload); err != nil {
		log.Println("Failed to unmarshal the JSON")
	}

	// pgx pool to connect to the db
	var isUsernameAvailable IsUsernameAvailable
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

	_, err = queries.GetUserFromUsername(ctx, usernameCheckPayload.Username)
	//error handling for this case
	if err != nil {
		// no rows found in the sql query
		if err == pgx.ErrNoRows {
			isUsernameAvailable.Success = true
			body, _ := json.Marshal(isUsernameAvailable)
			return events.APIGatewayProxyResponse{
				Body:       string(body),
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type":                 "application/json",
					"Access-Control-Allow-Origin":  "*", // Replace '*' with a specific domain for production
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
					"Access-Control-Allow-Origin":  "*", // Replace '*' with a specific domain for production
					"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type, Authorization",
				},
			}, nil
		}
	} else {
		// if a user was successfully found in this case
		isUsernameAvailable.Success = false
		body, _ := json.Marshal(isUsernameAvailable)
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*", // Replace '*' with a specific domain for production
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
