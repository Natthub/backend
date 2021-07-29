package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Database *mongo.Database

func ConnectDB() {
	credential := options.Credential{
		Username: "nat20344",
		Password: "0972256887",
	}

	clientOptions := options.Client().
		ApplyURI("mongodb+srv://nat20344:<0972256887>@cluster0.g1jof.mongodb.net/myFirstDatabase?retryWrites=true&w=majority").SetAuth(credential)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	Database = client.Database("linedb")
}
