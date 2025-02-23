package main

import (
	"context"
	"encoding/json"
	"fmt"

	"gc.yashk.dev/lambda/internal/db"
	"gc.yashk.dev/lambda/internal/middleware"
	"gc.yashk.dev/lambda/internal/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ListLessonsRequest struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Filter string `json:"filter"`
	Order  string `json:"order"`
}

type ListLessonsResponse struct {
	Lessons []types.Lesson `json:"lessons"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ctx := context.Background()
	var listLessonsRequest ListLessonsRequest
	if err := json.Unmarshal([]byte(request.Body), &listLessonsRequest); err != nil {
		fmt.Println("Failed to unmarshal", err.Error())
	}
	queries, pool, err := middleware.InitDb(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Something failed" + err.Error(),
		}, nil
	}
	defer pool.Close()

	var listLessonsResponse ListLessonsResponse

	res, err := queries.GetAllLessons(ctx, db.GetAllLessonsParams{
		IsPublic:   true,
		Limit:      int32(listLessonsRequest.Limit),
		Offset:     int32(listLessonsRequest.Offset),
		Similarity: listLessonsRequest.Filter,
		Column5:    listLessonsRequest.Order,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       err.Error(),
		}, nil
	}

	body, _ := json.Marshal(listLessonsResponse)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
