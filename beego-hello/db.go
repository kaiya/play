package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		cancelFunc()
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	quickstartDatabase := client.Database("quickstart")
	podcastCollection := quickstartDatabase.Collection("podcasts")
	episodesCollection := quickstartDatabase.Collection("episodes")

	//insert one
	podcastResult, err := podcastCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "My podcast"},
		{Key: "author", Value: "Kaiya"},
	})
	if err != nil {
		log.Fatal(err)
	}
	//insert many
	episodeResult, err := episodesCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "GraphQL for API Development"},
			{"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
			{"duration", 25},
		},
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "Progressive Web Application Development"},
			{"description", "Learn about PWA development with Tara Manicsic."},
			{"duration", 32},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %v documents into episode collection!\n", len(episodeResult.InsertedIDs))
}
