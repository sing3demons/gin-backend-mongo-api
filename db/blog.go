package db

import (
	"context"
	"time"

	"github.com/sing3demons/gin-backend-api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type blog struct {
	col *mongo.Collection
}

func NewBlog(col *mongo.Collection) *blog {
	return &blog{col}
}

func (tx blog) CreateBlog(body interface{}) (result *mongo.InsertOneResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return tx.col.InsertOne(ctx, body)
}

func (tx *blog) FindAll(filter interface{}) (results []models.Blog, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tx.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (tx *blog) FindById(filter interface{}) (result *models.Blog, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = tx.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
