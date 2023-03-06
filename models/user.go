package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Fullname  string             `json:"fullname" bson:"fullname,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Password  string             `json:"-" bson:"password,omitempty"`
	IsAdmin   bool               `json:"is_admin" bson:"is_admin,omitempty"`
	Blogs     []Blog             `json:"blogs" bson:"blogs,omitempty;default:[]"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	DeletedAt *time.Time         `json:"-" bson:"deleted_at,omitempty;default:null"`
}
