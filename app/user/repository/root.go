package UserRepository

import (
	"ocr.service.authorization/config"
	"ocr.service.authorization/model"
	"ocr.service.authorization/module/db"
)

type userRepository struct {
	db db.IDB
}

func (q *userRepository) Get(filter model.User) ([]model.User, error) {
	var arrUser []model.User
	err := q.db.Get(filter, &arrUser)
	return arrUser, err
}
func (q *userRepository) InsertOne(user model.User) (string, error) {
	return q.db.InsertOne(user)
}
func (q *userRepository) Update(filter model.User, user model.User) (int64, error) {
	return q.db.Update(filter, user)
}

func NewUserRepository() (*userRepository, error) {
	var q userRepository
	var err error
	CONFIG, err := config.NewConfig(nil)
	if err != nil {
		return nil, err
	}
	q.db, err = db.NewMongoRepository(CONFIG.GetString("MONGODB_DB"), "users")
	return &q, err
}
