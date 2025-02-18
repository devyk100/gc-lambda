package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"gc.yashk.dev/lambda/internal/middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ErrorResp struct {
	Message string `json:"message"`
}

type ConferenceAuthResponse struct {
	Ok bool `json:"Ok"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	contxt := context.TODO()

	conferenceAuth := ConferenceAuthResponse{
		Ok: false,
	}
	query, pool, err := middleware.InitDb(contxt)
	if err != nil {
		fmt.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error parsing form data: %v", err),
			StatusCode: 400,
		}, nil
	}
	defer pool.Close()

	values, err := url.ParseQuery(request.Body)
	if err != nil {
		fmt.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error parsing form data: %v", err),
			StatusCode: 403,
		}, nil
	}

	password := values.Get("password")
	email := values.Get("email")
	email += "@meet.yashk.dev"
	fmt.Println(password, email)
	res, err := query.GetLiveClassFromEmail(contxt, email)
	if err != nil {
		fmt.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error parsing start_time: %v", err),
			StatusCode: 403,
		}, nil
	}
	for _, class := range res {
		fmt.Println(strings.Split(class.Email, "@")[0])
		if class.Email == email && class.ModPassword == password {
			conferenceAuth.Ok = true
			jsonData, err := json.Marshal(conferenceAuth)
			if err != nil {
				fmt.Print(err.Error())
				return events.APIGatewayProxyResponse{
					Body:       fmt.Sprintf("Error marshaling JSON: %v", err),
					StatusCode: 500,
				}, nil
			}
			return events.APIGatewayProxyResponse{
				Body:       string(jsonData),
				StatusCode: 200,
			}, nil
		}
	}

	jsonData, err := json.Marshal(conferenceAuth)
	if err != nil {
		fmt.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error marshaling JSON: %v", err),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonData),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
