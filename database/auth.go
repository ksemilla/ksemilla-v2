package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/golang-jwt/jwt"
	"github.com/ksemilla/ksemilla-v2/config"
	"github.com/ksemilla/ksemilla-v2/graph/model"
)

func (db *DB) Login(input *model.LoginInput) (*model.LoginReturn, error) {

	config := config.Config()

	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result *bson.M
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&result)
	if err != nil {
		return nil, errors.New("cannot find user")
	}
	user_map := *result
	password := user_map["password"].(string)
	match := CheckPasswordHash(input.Password, password)
	_token_duration, _ := strconv.Atoi(os.Getenv("TOKEN_DURATION"))
	if match {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo":       "bar",
			"key":       "value",
			"ExpiresAt": int(time.Now().Add(time.Second * time.Duration(_token_duration)).Unix()),
			"userId":    user_map["_id"].(primitive.ObjectID).Hex(),
		})
		tokenString, _ := token.SignedString([]byte(config.APP_SECRET_KEY))

		jsonbody, _ := json.Marshal(user_map)
		user := model.User{}

		err = json.Unmarshal(jsonbody, &user)
		if err != nil {
			return nil, errors.New("json unmarshal error")
		}
		return &model.LoginReturn{
			User:  &user,
			Token: tokenString,
		}, nil
	} else {
		return nil, errors.New("wrong credentials")
	}
}

func (db *DB) VerifyToken(input *model.VerifyToken) (*model.User, error) {

	config := config.Config()

	token, err := jwt.Parse(input.Token, func(token *jwt.Token) (interface{}, error) {
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
			return nil, errors.New("expired token")
		}

		ObjectID, err := primitive.ObjectIDFromHex(claims["userId"].(string))
		if err != nil {
			log.Fatal(err)
		}
		collection := db.client.Database("ksemilla").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		user := model.User{}
		res := collection.FindOne(ctx, bson.M{"_id": ObjectID})
		err = res.Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("no user found")
			}
		}

		return &user, nil
	} else {
		fmt.Println(err)
		return nil, errors.New("token unrecognized")
	}
}

func (db *DB) ChangePassword(input *model.ChangePassword) (bool, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(input.ID)
	filter := bson.M{"_id": ObjectID}

	hash, _ := HashPassword(input.Password)

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "password", Value: hash},
	},
	}}
	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return false, errors.New("something went wrong")
	} else {
		return true, nil
	}
}
