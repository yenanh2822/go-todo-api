package users

import (
	"context"
	"fmt"
	"time"
	"todo_api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = utils.OpenCollection(utils.Client, "user")

func (user *User) Register() {
	// check exist username
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	count, err := userCollection.CountDocuments(ctx, bson.D{{Key: "username", Value: user.Username}})
	if count != 0 {
		msg := "This user already exists"
		fmt.Println(msg, err)
		// c.JSON(http.StatusConflict, gin.H{"msg": msg, "error": err})
		defer cancel()
		return
	}

	// insert
	_, err = userCollection.InsertOne(ctx, user)
	defer cancel()
	if err != nil {
		return
	}
}

func (user *User) Login() {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var checkUser User
	err := userCollection.FindOne(ctx, bson.D{{Key: "username", Value: user.Username}}).Decode(&checkUser)
	defer cancel()
	if err != nil {
		fmt.Println(err.Error())
		// c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not registered.", "error": err})
		return
	}
	if checkUser.Password != user.Password {
		fmt.Println(err.Error())
		// c.JSON(http.StatusUnauthorized, gin.H{"msg": "Incorrect password."})
		return
	}
}
