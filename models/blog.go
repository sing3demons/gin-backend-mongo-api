package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time         `json:"-" bson:"deleted_at,omitempty;default:null"`
	Topic     string             `bson:"topic" json:"topic"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
}
