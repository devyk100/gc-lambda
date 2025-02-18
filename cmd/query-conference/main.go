package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"gc.yashk.dev/lambda/internal/middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ErrorResp struct {
	Message string `json:"message"`
}

type Conference struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	MailOwner         string    `json:"mail_owner"`
	StartTime         time.Time `json:"start_time"`
	Duration          int       `json:"duration"`
	Password          string    `json:"password"`
	ModeratorPassword string    `json:"moderator_password"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	contxt := context.TODO()
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

	// consider name as the id here of the live class
	name := values.Get("name")
	startTimeStr := values.Get("start_time")
	mailOwner := values.Get("mail_owner")
	classId, err := strconv.Atoi(name)

	fmt.Println("The request", "name", name, "start_time", startTimeStr, "mail_owner", mailOwner)
	if err != nil {
		fmt.Print(err.Error(), "name is ", name)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error parsing start_time: %v", err),
			StatusCode: 403,
		}, nil
	}

	res, err := query.GetLiveClassFromId(contxt, int32(classId))
	if err != nil {
		fmt.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error parsing start_time: %v", err),
			StatusCode: 403,
		}, nil
	}

	// Parse the start_time string into a time.Time object
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		fmt.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error parsing start_time: %v", err),
			StatusCode: 403,
		}, nil
	}

	// if res.Email != mailOwner {

	// }
	if res.Email != mailOwner {
		errResp := ErrorResp{Message: "Other rooms not allowed"}
		jsonBody, err := json.Marshal(errResp)
		if err != nil {
			fmt.Print(err.Error())
		}
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBody),
			StatusCode: 200,
		}, nil
	}

	// Create a Conference object with the extracted data
	conference := Conference{
		ID:                int(res.ID), // You can generate a unique ID here if needed
		Name:              strconv.Itoa(classId),
		MailOwner:         mailOwner,
		StartTime:         startTime,
		Duration:          int(res.Length), // Set Duration to 1 year in milliseconds
		Password:          res.Password,
		ModeratorPassword: res.ModPassword,
	}

	// Convert the Conference object to JSON
	jsonData, err := json.Marshal(conference)
	if err != nil {
		fmt.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error marshaling JSON: %v", err),
			StatusCode: 500,
		}, nil
	}

	// Log the source IP (optional)
	sourceIP := request.RequestContext.Identity.SourceIP
	if sourceIP == "" {
		fmt.Println("Hello, world!")
	} else {
		fmt.Printf("Hello, %s!\n", sourceIP)
	}

	// Return the JSON response
	return events.APIGatewayProxyResponse{
		Body:       string(jsonData),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
