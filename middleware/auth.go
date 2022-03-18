package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ksemilla/ksemilla-v2/config"
	"github.com/ksemilla/ksemilla-v2/database"
	"github.com/ksemilla/ksemilla-v2/graph/model"
)

type key int

const (
	UserCtx key = iota
)

func GetUserCtx() key {
	return UserCtx
}

func UserContextBody(next http.Handler) http.Handler {

	var db = database.Connect()

	config := config.Config()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		userObj := model.User{}
		if len(authToken) == 0 {
			// IF NO TOKEN PROVIDED
			ctx := context.WithValue(r.Context(), UserCtx, nil)
			next.ServeHTTP(rw, r.WithContext(ctx))
		} else {
			token, _ := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					return nil, errors.New("unexpected signing method")
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(config.APP_SECRET_KEY), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				timeValue := int64(claims["ExpiresAt"].(float64)) - time.Now().Unix()

				if timeValue <= 0 {
					// return "", errors.New("expired token")
					fmt.Println("expired token")
					ctx := context.WithValue(r.Context(), UserCtx, nil)
					next.ServeHTTP(rw, r.WithContext(ctx))
				} else {
					user, err := db.FindOneUser(claims["userId"].(string))
					if err != nil {
						fmt.Println("auth:no user found with")
						ctx := context.WithValue(r.Context(), UserCtx, nil)
						next.ServeHTTP(rw, r.WithContext(ctx))
					} else {
						// THIS
						userObj = *user
						ctx := context.WithValue(r.Context(), UserCtx, userObj)
						next.ServeHTTP(rw, r.WithContext(ctx))
					}
				}
			} else {
				// return "", errors.New("token unrecognized")
				ctx := context.WithValue(r.Context(), UserCtx, nil)
				next.ServeHTTP(rw, r.WithContext(ctx))

			}
		}

		// js, _ := json.Marshal(&struct{ Email string }{"test@test.com"})
		// ctx := context.WithValue(r.Context(), "user", js)
		// next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
