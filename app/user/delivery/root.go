package UserDelivery

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	UserInterface "ocr.service.authorization/app/user/interface"
	UserUsecase "ocr.service.authorization/app/user/usecase"
	"ocr.service.authorization/config"
	"ocr.service.authorization/model"
)

type userDelivery struct {
	IdentityKey string
	useCase     UserInterface.IUserUseCase
}

// @tags User
// @Summary user
// @Description create user
// @start_time default
// @Param body body model.UserCreate true "json"
// @Success 200 {object} model.UserResponse	""
// @Router /api/v1/register [post]
func (q *userDelivery) Create(c *gin.Context) {
	var userCreate model.UserCreate
	var user model.User
	var userResponse model.UserResponse
	err := c.BindJSON(&userCreate)
	if err != nil {
		c.String(500, "can't parse body")
	}
	copier.Copy(&user, &userCreate)
	userResponse, err = q.useCase.Create(user)
	if err != nil {
		if err.Error() == "username already exists" {
			c.String(400, err.Error())
			return
		}
		c.String(500, err.Error())
		return
	}
	c.JSON(200, userResponse)
}

// @tags User
// @Summary user
// @Description create user
// @start_time default
// @Param user_id path string false "user id"
// @Param Authorization header string false "'Bearer ' + token"
// @Param body body model.UserUpdate true "json"
// @Success 200 {object} model.UserResponse	""
// @Router /api/v1/auth/user/{user_id} [post]
func (q *userDelivery) UpdateByID(c *gin.Context) {
	userId := c.Param("user_id")
	var userUpdate model.UserUpdate
	var user model.User
	if userId == "" {
		c.String(400, "require param user_id")
	}
	err := c.BindJSON(&userUpdate)
	if err != nil {
		c.String(500, "can't parse body")
	}
	copier.Copy(&user, &userUpdate)
	err = q.useCase.Update(userId, user)
	if err != nil {
		c.String(500, "insert to db failed")
	}
	c.Writer.WriteHeader(200)
}

func (q *userDelivery) Gets(c *gin.Context) {
	var user model.User
	err := c.BindQuery(&user)
	if err != nil {
		c.String(500, "can't parse body")
	}
	userResponse, err := q.useCase.GetByOwner(user)
	if err != nil {
		c.String(500, err.Error())
	}
	c.JSON(200, userResponse)
}

// @tags User
// @Summary user
// @Description create user
// @start_time default
// @Param user_id path string false "user id"
// @Param Authorization header string false "'Bearer ' + token"
// @Success 200 {object} model.UserResponse	""
// @Router /api/v1/auth/user/{user_id} [get]
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
