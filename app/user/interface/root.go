package UserInterface

import (
	"github.com/gin-gonic/gin"
	"ocr.service.authorization/model"
)

type IUserDelivery interface {
	Create(c *gin.Context)
	UpdateByID(c *gin.Context)
	Gets(c *gin.Context)
	GetByID(c *gin.Context)
}

type IUserUseCase interface {
	GetFull(user model.User) (model.User, error)
	GetByOwner(user model.User) (model.UserResponse, error)
	Create(user model.User) (model.UserResponse, error)
	Update(idUser string, user model.User) error
}

type IUserRepository interface {
	Get(filter model.User) ([]model.User, error)
	InsertOne(user model.User) (string, error)
	Update(filter model.User, user model.User) (int64, error)
}
