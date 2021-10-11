package models

import "time"

type User struct {
	Id           string    `json:"_id" bson:"_id,omitempty"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	First_name   string    `json:"first_name"`
	Last_name    string    `json:"last_name"`
	Created_date time.Time `json:"created_at" bson:"created_at,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `bson:"oldPassword" json:"oldPassword" validate:"nonnil,nonzero"`
	NewPassword string `bson:"newPassword" json:"newPassword" validate:"nonnil,nonzero"`
}
