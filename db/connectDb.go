package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://superadmin:123456@51.79.184.185:27017/"

func ConnectMongoDB() *mongo.Database {

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// user Connection database

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	userclient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {

	}

	// Check the connection
	err = userclient.Ping(ctx, nil)

	if err != nil {

	}

	fmt.Println("Connected to user MongoDB!")
	return userclient.Database("products")
}

//GetMongoDBClient , return mongo client for CRUD operations
/* func GetMongoDBClient() *mongo.Client {

    return userclient
} */
