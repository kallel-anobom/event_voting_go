package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kallel-anobom/event_voting_go/api/handler"
	"github.com/kallel-anobom/event_voting_go/api/repository"
	"github.com/kallel-anobom/event_voting_go/api/services/cache"
	"github.com/kallel-anobom/event_voting_go/api/services/database"
	"github.com/kallel-anobom/event_voting_go/api/services/pubsub"
	"github.com/kallel-anobom/event_voting_go/api/subscriber"
	"github.com/kallel-anobom/event_voting_go/api/usecase"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
)

func getRedisService() *cache.RedisService {
	redisAddr := "localhost:6379"
	redisPassword := "admin"
	redisDB := 0
	redisService, err := cache.NewRedisService(redisAddr, redisPassword, redisDB)
	if err != nil {
		fmt.Println(err)
		panic("Erro ao conectar ao Redis")
	}
	return redisService
}

func getMongoService() (*database.MongoService, error) {
	mongoURL := "mongodb://admin:admin@localhost:27017/?authSource=admin"
	return database.NewMongoService(mongoURL)
}

func getRabbitMQService() *pubsub.RabbitMQService {
	rabbitService, err := pubsub.NewRabbitMQService("amqp://admin:admin@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic("Erro ao conectar ao RabbitMQ")
	}
	return rabbitService
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	apiPort := "8001"

	mongoService, err := getMongoService()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mongoService.GetClient().Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	log.Println("MongoDB connection verified successfully!")

	testDB := mongoService.GetClient().Database("test_write")
	_, err = testDB.Collection("test").InsertOne(ctx, bson.M{"test": time.Now()})
	if err != nil {
		log.Fatalf("MongoDB write test failed: %v", err)
	}
	log.Println("MongoDB write test successful!")
	testDB.Collection("test").Drop(ctx)

	redisService := getRedisService()
	rabbitMQService := getRabbitMQService()

	votesRepository := repository.NewVotesRepository(mongoService)
	votesUsecase := usecase.NewVotesUsecase(votesRepository, redisService, rabbitMQService)
	votesHandler := handler.NewVotesHandler(votesUsecase)

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.POST("/api/vote", votesHandler.Vote)
	server.GET("/api/vote/summary", votesHandler.GetSummary)

	http.Handle("/metrics", promhttp.Handler())

	voteSubscriber := subscriber.NewVoteSubscribers(votesRepository, redisService, rabbitMQService)
	go voteSubscriber.SubscribeToPubsub()

	go http.ListenAndServe(":9090", nil)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server on port %s", apiPort)
		serverErr <- server.Run(":" + apiPort)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Fatalf("Server failed: %v", err)
	case <-sigs:
		log.Println("Shutting down gracefully...")

		// Fecha conexões na ordem inversa da criação
		rabbitMQService.Close()

		if err := redisService.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}

		if err := mongoService.Disconnect(); err != nil {
			log.Printf("Error closing MongoDB: %v", err)
		}

		log.Println("Server shutdown complete")
	}
}
