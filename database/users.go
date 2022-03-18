package database

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/ksemilla/ksemilla-v2/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) GetAllUsers() ([]*model.User, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.New("somethign went wrong with fetching users cursor")
	}
	var users []*model.User
	for cur.Next(ctx) {
		var user *model.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, errors.New("cannot decode a user")
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) IsCreateOwnerValid() (bool, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := model.CreateUser{}
	err := collection.FindOne(ctx, bson.M{}).Decode(&user)
	return err == mongo.ErrNoDocuments, nil
}

func (db *DB) CreateOwner(input *model.CreateOwner) (*model.User, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var OWNER_KEY = os.Getenv("OWNER_KEY")

	if OWNER_KEY != input.Key {
		return nil, errors.New("wrong key given")
	}

	user := model.CreateUser{}
	err := collection.FindOne(ctx, bson.M{}).Decode(&user)
	if err == mongo.ErrNoDocuments {

		user.Email = input.Email
		user.Role = "OWNER"

		hash, _ := HashPassword(input.Password)
		user.Password = hash

		res, err := collection.InsertOne(ctx, user)
		if err != nil {
			return nil, errors.New("error creating user")
		}

		return &model.User{
			ID:    res.InsertedID.(primitive.ObjectID).Hex(),
			Email: input.Email,
			Role:  "OWNER",
		}, nil
	}

	return nil, errors.New("cannot create starting user")
}

func (db *DB) FindOneUser(_id string) (*model.User, error) {
	ObjectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, err
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
}

func (db *DB) CreateUser(input *model.CreateUser) (*model.User, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": input.Email}
	val := model.User{}
	err := collection.FindOne(ctx, filter).Decode(&val)
	if err == mongo.ErrNoDocuments {
		if len(input.Password) > 0 {
			hash, _ := HashPassword(input.Password)
			input.Password = hash
		} else {
			hash, _ := HashPassword(RandStringRunes(6))
			input.Password = hash
		}

		res, err := collection.InsertOne(ctx, input)
		if err != nil {
			return nil, err
		}
		return &model.User{
			ID:    res.InsertedID.(primitive.ObjectID).Hex(),
			Email: input.Email,
			Role:  input.Role,
		}, nil
	} else {
		return nil, errors.New("existing email")
	}
}

func (db *DB) UpdateUser(input *model.UpdateUser) (*model.User, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(input.ID)
	filter := bson.M{"_id": ObjectID}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "role", Value: input.Role},
		{Key: "email", Value: input.Email},
	},
	}}
	res, err := collection.UpdateOne(ctx, filter, update)

	if res.ModifiedCount == 0 {
		return nil, errors.New("cant locate user")
	}

	if err != nil {
		return nil, err
	}

	user := model.User{}

	jsonbody, err := json.Marshal(*input)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonbody, &user); err != nil {
		return nil, err
	}
	user.ID = input.ID

	return &user, nil
}

func (db *DB) DeleteUser(id string) (*mongo.DeleteResult, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.New("cannot convert id to primitive type")
	}

	filter := bson.M{"_id": ObjectID}
	return collection.DeleteOne(ctx, filter)
}
