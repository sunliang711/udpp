package models

import (
	"context"
	"log"
	"time"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Mdb     *mongo.Client
	BlockDb *mongo.Client
)

func InitMongo(url string)  {
	var err error
	Mdb, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Mdb.Connect(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = Mdb.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("connect to mongodb error")
		panic(err)
	}
	log.Printf("connect to mongodb OK")
	return
}

func InitBlockDb(url string)  {
	var err error
	BlockDb, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = BlockDb.Connect(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = BlockDb.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("connect to mongodb error")
		panic(err)
	}
	log.Printf("connect to mongodb OK")
	return
}

