package db

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Database {
	var uri = viper.GetString("mongodb.connection")
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

	/* fmt.Println("Connected to user MongoDB!") */
	return userclient.Database(viper.GetString("mongodb.database"))
}

//GetMongoDBClient , return mongo client for CRUD operations
/* func GetMongoDBClient() *mongo.Client {

    return userclient
} */
