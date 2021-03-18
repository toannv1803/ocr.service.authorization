package UserUseCase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	UserInterface "ocr.service.authorization/app/user/interface"
	UserRepository "ocr.service.authorization/app/user/repository"
	"ocr.service.authorization/model"
	"ocr.service.authorization/module/salt"
	"time"
)

type userUseCase struct {
	repository UserInterface.IUserRepository
}

var nullUser = model.User{}
var nullUserResponse = model.UserResponse{}

func (q *userUseCase) GetByOwner(user model.User) (model.UserResponse, error) {
	var userResponse model.UserResponse
	user, err := q.GetFull(user)
	if err != nil {
		return nullUserResponse, err
	}
	err = copier.Copy(&userResponse, &user)
	if err != nil {
		return nullUserResponse, err
	}
	return userResponse, nil
}

func (q *userUseCase) GetFull(user model.User) (model.User, error) {
	arrUser, err := q.repository.Get(user)
	if err != nil {
		return nullUser, errors.New("get from db failed")
	}
	if len(arrUser) == 0 {
		return nullUser, errors.New("not found")
	}
	return arrUser[0], nil
}

func (q *userUseCase) Create(user model.User) (model.UserResponse, error) {
	_uuid := uuid.New().String()
	user.Id = _uuid
	user.Role = "user"
	user.Password = salt.HashAndSalt([]byte(user.Password))
	user.CreateAt = time.Now().Format(time.RFC3339)
	arrUser, err := q.repository.Get(model.User{Username: user.Username})
	if err != nil {
		return nullUserResponse, errors.New("get from db failed")
	}
	if len(arrUser) == 0 {
		_, err := q.repository.InsertOne(user)
		if err != nil {
			return nullUserResponse, errors.New("insert to db failed")
		}
	} else {
		return nullUserResponse, errors.New("username already exists")
	}
	var userResponse model.UserResponse
	err = copier.Copy(&userResponse, &user)
	if err != nil {
		return nullUserResponse, err
	}
	return userResponse, nil
}

func (q *userUseCase) Update(userId string, user model.User) error {
	user.Password = salt.HashAndSalt([]byte(user.Password))
	modifiedCount, err := q.repository.Update(model.User{Id: userId}, user)
	if err != nil {
		return errors.New("update user failed")
	}
	if modifiedCount == 0 {
		return errors.New("not found")
	}
	return nil
}

func NewUserUseCase() (UserInterface.IUserUseCase, error) {
	var q userUseCase
	var err error
	q.repository, err = UserRepository.NewUserRepository()
	if err != nil {
		return nil, err
	}
	return &q, err
}
