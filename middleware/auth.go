package gc_middleware

import (
	"context"

	"gc.yashk.dev/db"
	"gc.yashk.dev/env"
	"github.com/golang-jwt/jwt"
)

func JwtAuth(ctx context.Context, token *string, queries *db.Queries) (bool, error) {
	t, err := jwt.Parse(*token, func(token *jwt.Token) (interface{}, error) {
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
