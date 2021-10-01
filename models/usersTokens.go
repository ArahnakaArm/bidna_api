package models

import "time"

type User_token struct {
	Id           string    `json:"_id" bson:"_id,omitempty"`
	User_id      string    `json:"user_id"`
	Token        string    `json:"token"`
	Created_date time.Time `json:"created_at" bson:"created_at,omitempty"`
}
