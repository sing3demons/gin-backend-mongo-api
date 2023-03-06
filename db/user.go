package db

import (
	"context"
	"fmt"
	"time"

	"github.com/sing3demons/gin-backend-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	col *mongo.Collection
}

func New(col *mongo.Collection) *user {
	return &user{col}
}

func (u user) GetList(filter primitive.D) (results []models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lookupStage := bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "blogs"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "user_id"},
				{Key: "as", Value: "blogs"},
			},
		},
	}

	cur, err := u.col.Aggregate(ctx, mongo.Pipeline{lookupStage})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	fmt.Println(results)
	return results, nil
}

func (u user) Create(body interface{}) (result *mongo.InsertOneResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err = u.col.InsertOne(ctx, body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u user) CheckEmail(filter primitive.D) models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := models.User{}
	u.col.FindOne(ctx, filter).Decode(&result)

	return result
}

func (u *user) Search(filter primitive.D) (results *models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = u.col.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (u user) GetById(filter primitive.D) (result *models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = u.col.FindOne(ctx, filter).Decode(&result)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	return result, nil
}
