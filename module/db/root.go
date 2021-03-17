package db

type IDB interface {
	Get(filter interface{}, res interface{}) error
	InsertOne(data interface{}) (string, error)
	Update(filter interface{}, data interface{}) (int64, error)
}

func NewDB(dbName string, collectionName string) (IDB, error) {
	var q IDB
	var err error
	q, err = NewMongoRepository(dbName, collectionName)
	return q, err
}
