package users

import (
	"net/http"

	"todo_api/domain/users"
	"todo_api/service"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := service.UsersService.Register(user)
	// if err != nil {
	// 	msg := "User was not registered"
	// 	c.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "error": err})
	// 	return
	// }

	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := service.UsersService.Login(user)

	c.JSON(http.StatusOK, result)
}
