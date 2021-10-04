package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name" validate:"nonnil,nonzero"`
	Colors   []string           `bson:"colors" json:"colors" validate:"nonnil,nonzero"`
	ImageUrl string             `bson:"imageurl" json:"imageurl" validate:"nonnil,nonzero"`
	Category string             `bson:"category" json:"category" validate:"nonnil,nonzero"`
	Price    int                `bson:"price" json:"price" validate:"nonnil,nonzero"`
	Bidtype  string             `bson:"bidtype" json:"bidtype" validate:"nonnil,nonzero"`
}
