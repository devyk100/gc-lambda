package main

import (
	"context"
	"fmt"
	"log"

	db "gc.yashk.dev/checkusername/db_driver"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5/pgxpool"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dsn := ""
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Println(err.Error())
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer pool.Close()

	queries := db.New(pool)

	users, err := queries.GetUserFromUsername(ctx, "yash")
	if err != nil {
		fmt.Println(err.Error())
	}

	log.Print(users.Email, " ARE THE EMAILS OF THE USERS")

	var greeting string
	sourceIP := request.RequestContext.Identity.SourceIP

	if sourceIP == "" {
		greeting = "Hello, world from some other function!\n"
	} else {
		greeting = fmt.Sprintf("Hello, world from some other function!, %s!\n", sourceIP)
	}

	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func main() {

	lambda.Start(handler)
}
