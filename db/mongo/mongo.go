package mongo

import (
	"context"
	"time"

	"github.com/jasonlvhit/gocron"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Client
var isconnect = false
var errornum = 0

func Init(uri string) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	Mongo = client
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	isconnect = true
}

type mongolog struct {
	name string
}

func (m *mongolog) Write(p []byte) (n int, err error) {
	if !isconnect {
		return
	}
	bsonlog := make(bson.M)
	err = bson.UnmarshalExtJSON(p, true, &bsonlog)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = Mongo.Database(config.Config.Mongo.Datebase).Collection(m.name).InsertOne(ctx, bsonlog)
	if err != nil {
		errornum++
		if (errornum) > 10 {
			if isconnect {
				isconnect = false
				go reconnect()
			}
		}
	}
	return
}

var s = gocron.NewScheduler()

func reconnect() {
	s.Clear()
	_ = s.Every(1).Minutes().Do(task)
	<-s.Start()
}

func task() {
	if !isconnect {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		err := Mongo.Ping(ctx, nil)
		if err == nil {
			isconnect = true
			errornum = 0
			s.Remove(task)
			s.Clear()
		}
	}
}

func NewLog(name string) *mongolog {
	return &mongolog{
		name: name,
	}
}
