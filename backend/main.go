package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"live-polling-app/config"
	"live-polling-app/controllers"
	"live-polling-app/routes"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal("Invalid REDIS_URL:", err)
	}
	redisClient := redis.NewClient(redisOpts)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	db := mongoClient.Database(cfg.MongoDB)

	// Prevent the same email from being registered twice.
	_, _ = db.Collection("users").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
	)
	// One vote per browser/user id for each poll.
	_, _ = db.Collection("votes").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{Keys: bson.D{{Key: "pollId", Value: 1}, {Key: "voterId", Value: 1}}, Options: options.Index().SetUnique(true)},
	)

	authController := controllers.NewAuthController(db, cfg.JWTSecret)
	pollController := controllers.NewPollController(db, redisClient, cfg.JWTSecret)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://live-polling-app-3z5k.onrender.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))

	routes.Setup(r, authController, pollController, cfg.JWTSecret)

	port := cfg.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
