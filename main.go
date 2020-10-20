package main

import (
	"context"
	"fmt"
	"jwt-todo/auth-server/Config"
	"jwt-todo/auth-server/Models"
	"jwt-todo/auth-server/Routes"
	"log"
	"os"

	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var err error

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	dbConfig := Config.BuildConfig()
	dbURL := Config.DbURL(dbConfig)
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(dbURL)
	Config.CLI, err = mongo.Connect(ctx, clientOpts)

	if err != nil {
		fmt.Println("statuse: ", err)
	}

	dsn := os.Getenv("REDIS_DSN")
	fmt.Println(dsn)
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	Models.Client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := Models.Client.Ping().Result()
	if err != nil {
		panic(err)
	}

	router := Routes.SetupRouter()

	log.Fatal(router.Run(":" + port))
}
