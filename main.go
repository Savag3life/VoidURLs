package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortURL struct {
	DatabaseId primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ID         string             `json:"id" bson:"id"`
	OriginURL  string             `json:"origin_url" bson:"origin_url"`
	Views      int                `json:"views" bson:"views"`
}

type RequestURL struct {
	InputURL string `json:"url"`
}

func main() {
	enverr := godotenv.Load(".env")

	if enverr != nil {
		log.Fatal("Failed to load .env file")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}
	defer client.Disconnect(ctx)

	urlscollection := client.Database("urls").Collection("urls")

	router.GET("/:value", func(c *gin.Context) {
		value := c.Param("value")
		urlscollection := client.Database("urls").Collection("urls")
		filter := bson.D{{Key: "id", Value: value}}

		var result ShortURL
		urlscollection.FindOne(ctx, filter).Decode(&result)

		updatefilter := bson.M{"_id": result.DatabaseId}
		update := bson.M{
			"$inc": bson.M{
				"views": 1,
			},
		}
		urlscollection.UpdateOne(ctx, updatefilter, update, options.Update().SetUpsert(false))
		c.Redirect(http.StatusMovedPermanently, result.OriginURL)
	})

	router.POST("/void/create", func(c *gin.Context) {
		var url RequestURL
		if err := c.Bind(&url); err != nil {
			c.JSON(503, gin.H{"error": "Internal Server Error - Database Unavailable"})
		}

		err, id := generateUniqueRandom(urlscollection, ctx)

		if err != nil {
			c.JSON(503, gin.H{"error": err.Error()})
		}

		fmt.Printf("Creating short URL for %s pointing to %s\n", url.InputURL, id)
		urlscollection.InsertOne(ctx, ShortURL{DatabaseId: primitive.NewObjectID(), ID: id, OriginURL: url.InputURL, Views: 0})
		c.JSON(201, id)
	})

	router.Run(os.Getenv("ADDRESS") + ":" + os.Getenv("PORT"))
}

func generateUniqueRandom(urlscollection *mongo.Collection, ctx context.Context) (err error, random string) {

	random = string(generateRandom())
	filter := bson.D{{Key: "id", Value: random}}
	for value := urlscollection.FindOne(ctx, filter); value.Err() == nil; {
		random = string(generateRandom())
	}

	return nil, random

}

func generateRandom() []byte {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return result
}
