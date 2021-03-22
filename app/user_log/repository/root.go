package UserLogRepository

import (
	"ocr.service.authorization/config"
	"ocr.service.authorization/model"
	"ocr.service.authorization/module/db"
)

type userLogRepository struct {
	db db.IDB
}

func (q *userLogRepository) Get(filter model.UserLog) ([]model.UserLog, error) {
	var arrUserLog []model.UserLog
	err := q.db.Get(filter, &arrUserLog)
	return arrUserLog, err
}
func (q *userLogRepository) InsertOne(userLog model.UserLog) (string, error) {
	return q.db.InsertOne(userLog)
}
func (q *userLogRepository) Update(filter model.UserLog, userLog model.UserLog) (int64, error) {
	return q.db.Update(filter, userLog)
}

func NewUserLogRepository() (*userLogRepository, error) {
	var q userLogRepository
	var err error
	CONFIG, err := config.NewConfig(nil)
	if err != nil {
		return nil, err
	}
	q.db, err = db.NewMongoRepository(CONFIG.GetString("MONGODB_DB"), "user_logs")
	return &q, err
}
