package main

import (
	"context"
	"encoding/json"
	"fmt"

	"gc.yashk.dev/lambda/internal/middleware"
	"gc.yashk.dev/lambda/internal/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type GetCoursesRequest struct {
	Id int `json:"id"`
}

type GetCourseResponse struct {
	Course types.Course `json:"course"`
	Uid    int          `json:"uid"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ctx := context.Background()
	var getCoursesRequest GetCoursesRequest
	if err := json.Unmarshal([]byte(request.Body), &getCoursesRequest); err != nil {
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

	var getCourseResponse GetCourseResponse

	res, err := queries.GetCourse(ctx, int32(getCoursesRequest.Id))

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       err.Error(),
		}, nil
	}

	getCourseResponse.Course.Description = res.Description
	getCourseResponse.Course.Id = int(res.ID)
	getCourseResponse.Course.ImageUrl = res.ImgUrl
	getCourseResponse.Course.Instructor = res.UserName
	getCourseResponse.Course.InstructorUrl = res.UserPicture
	getCourseResponse.Course.Language = res.Language
	getCourseResponse.Course.Name = res.Name
	getCourseResponse.Uid = int(res.UserID)

	body, _ := json.Marshal(getCourseResponse)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
