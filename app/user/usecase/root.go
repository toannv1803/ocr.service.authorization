package UserUseCase

import (
	"errors"
	"github.com/google/uuid"
	UserInterface "ocr.service.authorization/app/user/interface"
	UserRepository "ocr.service.authorization/app/user/repository"
	"ocr.service.authorization/model"
	"ocr.service.authorization/module/salt"
)

type userUseCase struct {
	repository UserInterface.IUserRepository
}

var nullUser = model.User{}

func (q *userUseCase) GetByOwner(user model.User) (model.User, error) {
	user, err := q.GetFull(user)
	if err != nil {
		return nullUser, err
	}
	user.Password = ""
	return user, nil
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

func (q *userUseCase) Create(user model.User) (model.User, error) {
	_uuid := uuid.New().String()
	user.Id = _uuid
	user.Roles = "user"
	user.Password = salt.HashAndSalt([]byte(user.Password))
	arrUser, err := q.repository.Get(model.User{Username: user.Username})
	if err != nil {
		return nullUser, errors.New("get from db failed")
	}
	if len(arrUser) == 0 {
		_, err := q.repository.InsertOne(user)
		if err != nil {
			return nullUser, errors.New("insert to db failed")
		}
	} else {
		return nullUser, errors.New("username already exists")
	}
	user.Password = ""
	return user, nil
}

func (q *userUseCase) Update(userId string, user model.User) error {
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
