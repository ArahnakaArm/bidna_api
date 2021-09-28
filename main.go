package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gofiber/routes"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://superadmin:123456@51.79.184.185:27017/"

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	/* 	// Ping the primary
	   	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
	   		panic(err)
	   	} */
	fmt.Println("Successfully connected and pinged.")

	db := client.Database("products")

	collection := db.Collection("products")
	/* 	filter := bson.M{"name": "test1"}
	   	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	   	cur ,err := collection.Find(ctx, bson.M{}).Decode(&result)
	   	if err != nil {

	   	}
	*/
	/* 	cursor, err := collection.Find(ctx, bson.M{})
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	defer cursor.Close(ctx)
	   	for cursor.Next(ctx) {
	   		var episode bson.M
	   		if err = cursor.Decode(&episode); err != nil {
	   			log.Fatal(err)
	   		}
	   		fmt.Println(episode)
	   	} */

	var episodes []Result
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(ctx, &episodes); err != nil {
		panic(err.Error())
	}

	fmt.Println(episodes)

	b, err := json.Marshal(episodes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	/* 	print(result.Name)
	 */
	/* findCursor, findErr := collection.Find(context.TODO(), bson.D{})
	if findErr != nil {
		panic(findErr)
	}
	var findResults []bson.M
	if findErr = findCursor.All(context.TODO(), &findResults); findErr != nil {
		panic(findErr)
	}
	for _, result := range findResults {
		fmt.Println(result)
	} */

	routes.AddCommonRoute()

	/* 	http.HandleFunc("/apix/", func(w http.ResponseWriter, r *http.Request) {
	   		fmt.Fprintf(w, "Hello World")
	   	})
	*/
	http.HandleFunc("/apix/greet/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/greet/"):]
		fmt.Fprintf(w, "Hello %s\n", name)
	})
	http.ListenAndServe("127.0.0.1:3334", nil)
	/* http.ListenAndServe(":8000", nil) */
}

type Result struct {
	Name   string   `bson:"name,omitempty"`
	Colors []string `bson:"colors,omitempty"`
}
