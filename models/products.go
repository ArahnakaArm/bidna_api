package models

type Product struct {
	Name   string   `bson:"name,omitempty"`
	Colors []string `bson:"colors,omitempty"`
}
