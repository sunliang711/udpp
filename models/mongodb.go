package models

import (
	"context"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/sirupsen/logrus"
)

var (
	Mdb     *mongo.Client
	BlockDb *mongo.Client
)

func InitMongo(url string) {
	var err error
	Mdb, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = Mdb.Connect(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = Mdb.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.Fatalf("Connect to mongodb: %v failed", url)
	}
	logrus.Infof("Connected to mongodb: %v", url)
	return
}

func InitBlockDb(url string) {
	var err error
	BlockDb, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = BlockDb.Connect(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = BlockDb.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.Fatalf("Connect to blockDB: %v failed", url)
		panic(err)
	}
	logrus.Infof("Connected to blockDB: %v", url)
	return
}
