package database

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ksemilla/ksemilla-v2/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) CreateInvoice(input *model.CreateInvoice) *model.Invoice {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
		// return &model.Invoice{}
	}
	return &model.Invoice{
		ID:          res.InsertedID.(primitive.ObjectID).Hex(),
		From:        input.From,
		DateCreated: input.DateCreated,
		Address:     input.Address,
		Amount:      input.Amount,
	}
}

func (db *DB) GetInvoice(id string) (*model.Invoice, error) {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ObjectID}
	res := collection.FindOne(ctx, filter)
	invoice := model.Invoice{}
	err := res.Decode(&invoice)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no document found")
	}
	return &invoice, nil
}

func (db *DB) GetAllInvoices(page int64) (*model.PaginatedInvoices, error) {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	total, _ := collection.CountDocuments(ctx, bson.M{})
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "datecreated", Value: -1}})
	var perPage int64 = 5
	findOptions.SetSkip((page - 1) * perPage)
	findOptions.SetLimit(perPage)

	cur, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var invoices []*model.Invoice
	for cur.Next(ctx) {
		var invoice *model.Invoice
		err := cur.Decode(&invoice)
		if err != nil {
			log.Fatal(err)
		}
		invoices = append(invoices, invoice)
	}
	return &model.PaginatedInvoices{
		Data:  invoices,
		Total: int(total),
	}, nil
}

func (db *DB) UpdateInvoice(input *model.UpdateInvoice) (*model.Invoice, error) {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(input.ID)
	filter := bson.M{"_id": ObjectID}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "address", Value: input.Address},
		{Key: "from", Value: input.From},
		{Key: "dateCreated", Value: input.DateCreated},
		{Key: "amount", Value: input.Amount},
	},
	}}
	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	invoice := model.Invoice{}
	// res.Decode(&invoice)

	jsonbody, err := json.Marshal(*input)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonbody, &invoice); err != nil {
		return nil, err
	}
	invoice.ID = input.ID

	return &invoice, nil
}

func (db *DB) DeleteInvoice(id string) (*mongo.DeleteResult, error) {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.New("cannot process id")
	}

	filter := bson.M{"_id": ObjectID}
	return collection.DeleteOne(ctx, filter)
}

func (db *DB) Test() (*model.Invoice, error) {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "DateCreated", Value: 6},
	},
	}}
	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
