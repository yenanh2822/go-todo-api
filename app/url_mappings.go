package app

import (
	"net/http"

	"todo_api/controllers/users"
	services "todo_api/service"

	middlewares "todo_api/middleware"

	"github.com/gin-gonic/gin"
)

func mapUrls() {
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
	router.POST("/register", users.Register)
	// sign in endpoint
	router.POST("/signin", users.Login)
	router.Group("/task").Use(middlewares.JwtAuthMiddleware())
}
