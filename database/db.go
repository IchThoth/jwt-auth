package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = createDB()

func createDB() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoUrl := os.Getenv("MONGO_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created mongodb connection")

	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("cluster0").Collection(collectionName)
	return collection
}
