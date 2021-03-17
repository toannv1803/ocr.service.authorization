package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"ocr.service.authorization/config"
	"os"
	"time"
)

type mongoRepository struct {
	collection *mongo.Collection
}

func (q *mongoRepository) Get(filter interface{}, res interface{}) error {
	//var _filter interface{} = filter
	//if filter.Id != "" {
	//	objectId, err := primitive.ObjectIDFromHex(filter.Id)
	//	if err != nil {
	//		return nil, err
	//	}
	//	_filter = bson.M{"_id": objectId}
	//}
	//doc, err := toDoc(filter)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := q.collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	//var arrOriginModel []bson.M
	if err = cur.All(ctx, res); err != nil {
		return err
	}
	defer cur.Close(ctx)
	return nil
}
func (q *mongoRepository) InsertOne(data interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := q.collection.InsertOne(ctx, data)
	if err == nil {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
	return "", err
}
func (q *mongoRepository) Delete(image interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := q.collection.DeleteMany(ctx, image)
	if err == nil {
		return result.DeletedCount, err
	}
	return 0, err
}
func (q *mongoRepository) Update(filter interface{}, obj interface{}) (int64, error) {
	//var _filter interface{} = filter
	//if filter.Id != "" {
	//	objectId, err := primitive.ObjectIDFromHex(filter.Id)
	//	if err != nil {
	//		return err
	//	}
	//	_filter = bson.M{"_id": objectId}
	//}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := q.collection.UpdateOne(ctx, filter, bson.D{{"$set", obj}})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, err
}

func NewMongoRepository(dbName string, collectionName string) (*mongoRepository, error) {
	var p mongoRepository
	client, err := newMongoClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	collection := client.Database(dbName).Collection(collectionName)
	p.collection = collection
	return &p, nil
}

func newMongoClient() (*mongo.Client, error) {
	CONFIG, err := config.NewConfig(nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//return  mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	var client *mongo.Client
	if CONFIG.GetString("MONGODB_USERNAME") != "" {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+CONFIG.GetString("MONGODB_USERNAME")+":"+CONFIG.GetString("MONGODB_PASSWORD")+"@"+CONFIG.GetString("MONGODB_HOST")+":"+CONFIG.GetString("MONGODB_PORT")+"/?authSource="+CONFIG.GetString("MONGODB_DB")))
	} else {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+CONFIG.GetString("MONGODB_HOST")+":"+CONFIG.GetString("MONGODB_PORT")+"/?authSource="+CONFIG.GetString("MONGODB_DB")))
	}
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
