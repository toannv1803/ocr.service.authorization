package UserDelivery

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	UserInterface "ocr.service.authorization/app/user/interface"
	UserUsecase "ocr.service.authorization/app/user/usecase"
	"ocr.service.authorization/config"
	"ocr.service.authorization/model"
)

type userDelivery struct {
	IdentityKey string
	useCase     UserInterface.IUserUseCase
}

func (q *userDelivery) Create(c *gin.Context) {
	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		c.String(500, "can't parse body")
	}
	user, err = q.useCase.Create(user)
	if err != nil {
		if err.Error() == "username already exists" {
			c.String(400, err.Error())
			return
		}
		c.String(500, err.Error())
		return
	}
	c.JSON(200, user)
}

func (q *userDelivery) UpdateByID(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.String(400, "require param user_id")
	}
	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		c.String(500, "can't parse body")
	}
	err = q.useCase.Update(userId, user)
	if err != nil {
		c.String(500, "insert to db failed")
	}
}

func (q *userDelivery) Gets(c *gin.Context) {
	var user model.User
	err := c.BindQuery(&user)
	if err != nil {
		c.String(500, "can't parse body")
	}
	user, err = q.useCase.GetByOwner(user)
	if err != nil {
		c.String(500, err.Error())
	}
	c.JSON(200, user)
}

func (q *userDelivery) GetByID(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.String(400, "require param user_id")
	}
	claims := jwt.ExtractClaims(c)
	ownerID := claims[q.IdentityKey]
	if ownerID == userId {
		user, err := q.useCase.GetByOwner(model.User{Id: userId})
		if err != nil {
			c.String(500, err.Error())
		}
		c.JSON(200, user)
	} else {
		c.String(401, "not allow")
	}
}

func NewUserDelivery() (UserInterface.IUserDelivery, error) {
	CONFIG, _ := config.NewConfig(nil)
	var q userDelivery
	var err error
	q.useCase, err = UserUsecase.NewUserUseCase()
	q.IdentityKey = CONFIG.GetString("IDENTITY_KEY")
	return &q, err
}
