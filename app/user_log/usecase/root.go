package UserLogUseCase

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	UserLogInterface "ocr.service.authorization/app/user_log/interface"
	"ocr.service.authorization/config"
	"ocr.service.authorization/enum"
	"ocr.service.authorization/model"
	"ocr.service.authorization/module/db"
	"time"
)

type userLogUseCase struct {
	db                 db.IDB
	maxConcurrentToken int
	expiredTime        time.Duration
}

func (q *userLogUseCase) Add(userLog model.UserLog) error {
	// get all token has created 1h ago
	// count unexpired token
	// check limit
	// add user_log
	_, err := q.db.InsertOne(userLog)
	return err
}

func (q *userLogUseCase) IsAllowLogin(userId string) (bool, error) {
	arrUserLog, err := q.Gets(userId, time.Now().Add(-q.expiredTime))
	if err != nil {
		return false, err
	}
	countUnexpiredToken := 0
	for i := range arrUserLog {
		expiredTime, err := time.Parse(time.RFC3339, arrUserLog[i].ExpiredTime)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expiredTime) && arrUserLog[i].Status == enum.UserLogStatusInit { // chua het han
			countUnexpiredToken++
		}
	}
	if countUnexpiredToken >= q.maxConcurrentToken {
		return false, nil
	}
	return true, nil
}

func (q *userLogUseCase) Gets(userId string, startTime time.Time) ([]model.UserLog, error) {
	var arrUserLog []model.UserLog
	err := q.db.Get(bson.M{
		"user_id": userId,
		"expired_time": bson.M{
			"$gt": startTime.Format(time.RFC3339),
			//"$lt": toDate,
		},
	}, &arrUserLog)
	return arrUserLog, err
}

func (q *userLogUseCase) Update(filter model.UserLog, data model.UserLog) (int64, error) {
	nModify, err := q.db.Update(filter, data)
	if err != nil {
		return 0, errors.New("update user_log failed")
	}
	return nModify, err
}

func NewUserLogUseCase() (UserLogInterface.IUserLogUseCase, error) {
	var q userLogUseCase
	var err error
	CONFIG, err := config.NewConfig(nil)
	q.maxConcurrentToken = CONFIG.GetInt("LIMIT_CONCURRENT_LOGIN")
	q.expiredTime = CONFIG.GetDuration("TOKEN_EXPIRE_TIME")
	q.db, err = db.NewMongoRepository(CONFIG.GetString("MONGODB_DB"), "user_logs")
	if err != nil {
		return nil, err
	}
	return &q, err
}
