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

type ListAllCoursesRequest struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Filter string `json:"filter"`
	Order  string `json:"order"`
}

type ListAllCourseResponse struct {
	Courses []types.Course `json:"courses"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ctx := context.Background()
	var listAllCoursesRequest ListAllCoursesRequest
	if err := json.Unmarshal([]byte(request.Body), &listAllCoursesRequest); err != nil {
		fmt.Println("Failed to unmarshal", err.Error())
	}
	queries, pool, err := middleware.InitDb(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "",
		}, nil
	}
	defer pool.Close()

	var listAllCourseResponse ListAllCourseResponse
	res, err := queries.GetAllCourses(ctx, db.GetAllCoursesParams{
		Limit:      int32(listAllCoursesRequest.Limit),
		Offset:     int32(listAllCoursesRequest.Offset),
		IsPublic:   true,
		Similarity: listAllCoursesRequest.Filter,
		Column5:    listAllCoursesRequest.Order,
	})
	if err != nil {
		fmt.Println("Some error occured", err.Error())
	}

	for _, val := range res {
		listAllCourseResponse.Courses = append(listAllCourseResponse.Courses, types.Course{
			Description:   val.Description,
			Id:            int(val.ID),
			Language:      val.Language,
			Name:          val.Name,
			Instructor:    val.UserName,
			ImageUrl:      val.ImgUrl,
			InstructorUrl: val.Picture,
		})
	}
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "",
		}, nil
	}

	body, _ := json.Marshal(listAllCourseResponse)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
