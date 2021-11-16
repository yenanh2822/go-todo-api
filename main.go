package main

import (
	"net/http"
	"os"
	database "todo_api/service"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "ahihi"})
	})
	// get all tasks endpoint
	router.GET("/task/:id", database.GetTask)
	// get task by id endpoint
	router.GET("/task", database.GetTasks)
	// create task endpoint
	router.POST("/task", database.CreateTask)
	// update task endpoint
	router.PUT("/task/:id", database.UpdateTask)
	// delete task endpoint
	router.DELETE("/task/:id", database.DeleteTask)
	router.Run(":" + port)
}
