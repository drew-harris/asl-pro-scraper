package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func GetDatabase() *mongo.Database {
	return GetConnection().Database(os.Getenv("MONGO_DATABASE"))
}

func init() {
	godotenv.Load(".env")
	fmt.Print("Connecting to MongoDB Atlas...")
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		panic("Empty MongoDB URI")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println("Connection Successful")
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println("Atlas Ping Successful")

}

func GetConnection() *mongo.Client {
	return client
}

func CloseDBConnection() {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")
}
