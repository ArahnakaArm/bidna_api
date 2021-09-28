package models

type Product struct {
	Name   string   `bson:"name,omitempty" json:"name"`
	Colors []string `bson:"colors,omitempty" json:"colors"`
}
