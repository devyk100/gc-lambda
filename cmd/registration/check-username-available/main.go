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

type IsUsernameAvailable struct {
	Success bool `json:"success"`
}

type UsernameCheckPayload struct {
	Username string `json:"username"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ctx := context.Background()
	// unmarshalling the hjson from the request
	var usernameCheckPayload UsernameCheckPayload
	if err := json.Unmarshal([]byte(request.Body), &usernameCheckPayload); err != nil {
		log.Println("Failed to unmarshal the JSON")
	}

	var isUsernameAvailable IsUsernameAvailable

	// init the DB and get the queries client
	queries, pool, err := middleware.InitDb(ctx)
	if err != nil {
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
	defer pool.Close()

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
		isUsernameAvailable.Success = false
		body, _ := json.Marshal(isUsernameAvailable)
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
