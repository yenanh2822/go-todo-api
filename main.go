package main

import (
	"net/http"
	"os"
	services "todo_api/service"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "ahihi"})
	})
	// get all tasks endpoint
	router.GET("/task/:id", services.GetTask)
	// get task by id endpoint
	router.GET("/task", services.GetTasks)
	// create task endpoint
	router.POST("/task", services.CreateTask)
	// update task endpoint
	router.PUT("/task/:id", services.UpdateTask)
	// delete task endpoint
	router.DELETE("/task/:id", services.DeleteTask)
	// register endpoint
	router.POST("/register", services.Register)
	// signin endpoint
	router.POST("/signin", services.Login)
	router.Run(":" + port)
}
