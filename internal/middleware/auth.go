package middleware

import (
	"context"
	"fmt"
	"strings"

	"gc.yashk.dev/lambda/internal/db"
	"gc.yashk.dev/lambda/internal/env"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

func JwtAuth(ctx context.Context, request *events.APIGatewayProxyRequest, queries *db.Queries) (bool, error) {
	authHeader, exists := request.Headers["Authorization"]
	if !exists {
		return false, nil
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		fmt.Println("Invalid Authorization header format")
		return false, nil
	}

	bearerType := parts[0]
	fmt.Println("Bearer Type:", bearerType)

	token := parts[1]

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(env.NEXTAUTH_SECRET)
		return secret, nil
	})
	if err != nil {
		return false, err
	}
	if !t.Valid {
		return false, nil
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	email, emailExists := claims["email"].(string)
	if !emailExists {
		return false, nil
	}

	user, err := queries.GetUserFromEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if user.Email == email {
		return true, nil
	} else {
		return false, nil
	}
}
