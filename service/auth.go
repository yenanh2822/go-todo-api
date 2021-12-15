package services

import (
	"context"
	"net/http"
	"time"
	utils "todo_api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//connect to to the database and open a user collection
var userCollection *mongo.Collection = utils.OpenCollection(utils.Client, "user")

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	User_id  string             `json:"User_id"`
	Username string             `json:"username" binding:"required"`
	Password string             `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	// check exist username
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	count, err := userCollection.CountDocuments(ctx, bson.D{{Key: "username", Value: user.Username}})
	if count != 0 {
		msg := "This user already exists"
		c.JSON(http.StatusConflict, gin.H{"msg": msg, "error": err})
		defer cancel()
		return
	}

	// insert
	result, err := userCollection.InsertOne(ctx, user)
	defer cancel()
	if err != nil {
		msg := "User was not registered"
		c.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "error": err})
		return
	}

	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var checkUser User
	err := userCollection.FindOne(ctx, bson.D{{Key: "username", Value: user.Username}}).Decode(&checkUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not registered.", "error": err})
		return
	}
	if checkUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Incorrect password."})
		return
	}
	token, err := utils.GenerateToken(checkUser.User_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot generate token.", "error": err})
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
