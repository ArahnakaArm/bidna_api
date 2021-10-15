package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"nonnil,nonzero"`
	Colors    []string           `bson:"colors" json:"colors" validate:"nonnil,nonzero"`
	ImageUrl  string             `bson:"imageUrl" json:"imageUrl" validate:"nonnil,nonzero"`
	Category  string             `bson:"category" json:"category" validate:"nonnil,nonzero"`
	Price     int                `bson:"price" json:"price" validate:"nonnil,nonzero"`
	Bidtype   string             `bson:"bidtype" json:"bidtype" validate:"nonnil,nonzero"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"deletedAt"`
}
