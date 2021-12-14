package database

import (
	"context"
	"net/http"
	"time"
	token "todo_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// create a validator object
var validate = validator.New()

//connect to to the database and open a task collection
var taskCollection *mongo.Collection = OpenCollection(Client, "task")

type Task struct {
	ID          primitive.ObjectID `bson:"_id"`
	Task_id     string             `json:"task_id"`
	Name        string             `json:"name" validate:"required,min=2,max=100"`
	Description string             `json:"description" validate:"max=1000"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	User_id     string             `json:"user_id"`
}

func CreateTask(c *gin.Context) {
	user_id := getUserId(c)

	var task Task
	// validate
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New task is empty"})
		return
	}
	validationErr := validate.Struct(task)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// assign value
	task.Created_at = time.Now()
	task.Updated_at = time.Now()
	task.ID = primitive.NewObjectID()
	task.Task_id = task.ID.Hex()
	task.User_id = user_id

	// insert
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := taskCollection.InsertOne(ctx, task)
	defer cancel()

	if err != nil {
		msg := "Task was not created"
		c.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func UpdateTask(c *gin.Context) {
	user_id := getUserId(c)
	var task, curTask Task
	var id string = c.Param("id")
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is empty"})
		return
	}
	task.Updated_at = time.Now()

	// find task
	filter := bson.D{{Key: "task_id", Value: id}, {Key: "user_id", Value: user_id}}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	err := taskCollection.FindOne(ctx, filter).Decode(&curTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No matched task"})
			defer cancel()
			return
		}
	}
	// update task
	if task.Name == "" {
		task.Name = curTask.Name
	}
	if task.Description == "" {
		task.Description = curTask.Description
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: task.Name}, {Key: "description", Value: task.Description}, {Key: "updated_at", Value: task.Updated_at}}}}
	result, err := taskCollection.UpdateOne(ctx, filter, update)
	defer cancel()
	if err != nil {
		msg := "Task was not updated"
		c.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "error": err})
		return
	}

	c.JSON(http.StatusOK, result)
}

func DeleteTask(c *gin.Context) {
	user_id := getUserId(c)
	var id string = c.Param("id")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	filter := bson.D{{Key: "task_id", Value: id}, {Key: "user_id", Value: user_id}}
	result, err := taskCollection.DeleteOne(ctx, filter)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No matched task"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetTasks(c *gin.Context) {
	user_id := getUserId(c)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := taskCollection.Find(ctx, bson.D{{Key: "user_id", Value: user_id}})
	defer cancel()
	if err != nil {
		msg := "Cannot get tasks"
		c.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "error": err})
		return
	} else {
		for result.Next(ctx) {
			var trueResult bson.M
			result.Decode(&trueResult)
			c.JSON(http.StatusOK, trueResult)
		}
	}
}

func GetTask(c *gin.Context) {
	user_id := getUserId(c)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var id string = c.Param("id")
	var result bson.M
	err := taskCollection.FindOne(ctx, bson.D{{Key: "task_id", Value: id}, {Key: "user_id", Value: user_id}}).Decode(&result)
	defer cancel()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No matched task"})
			return
		}
	}
	c.JSON(http.StatusOK, result)
}

func getUserId(c *gin.Context) string {
	user_id, errs := token.ExtractTokenID(c)
	if errs != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errs.Error()})
		return ""
	}
	return user_id
}
