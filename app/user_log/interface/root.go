package UserLogInterface

import (
	"ocr.service.authorization/model"
	"time"
)

type IUserLogUseCase interface {
	Add(userLog model.UserLog) error
	IsAllowLogin(userId string) (bool, error)
	Gets(userId string, startTime time.Time) ([]model.UserLog, error)
	Update(filter model.UserLog, data model.UserLog) (int64, error)
}
