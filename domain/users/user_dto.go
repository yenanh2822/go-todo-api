package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	User_id  string             `json:"User_id"`
	Username string             `json:"username" binding:"required"`
	Password string             `json:"password" binding:"required"`
}
