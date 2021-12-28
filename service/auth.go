package service

import (
	"fmt"
	"todo_api/domain/users"
	utils "todo_api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//connect to to the database and open a user collection
var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	Register(users.User) *users.User
	Login(users.User) string
}

func (s *usersService) Register(user users.User) *users.User {

	validationErr := validate.Struct(user)
	if validationErr != nil {
		fmt.Println(validationErr)
		// c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})

		// tupt4: should return err
		return nil
	}
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	user.Register()

	return &user
}

func (s *usersService) Login(user users.User) string {

	validationErr := validate.Struct(user)
	if validationErr != nil {
		fmt.Println(validationErr)

		// tupt4: should return err
		return "error"
	}

	user.Login()

	token, err := utils.GenerateToken(user.User_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	return token
}
