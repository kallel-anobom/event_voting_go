package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kallel-anobom/event_voting_go/controller"
	"github.com/kallel-anobom/event_voting_go/repository"
	"github.com/kallel-anobom/event_voting_go/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo .env - %v", err)
		log.Println("Continuando com variáveis de ambiente do sistema")
	}
	
	gin.SetMode(gin.ReleaseMode)
	server := gin.New() 
	server.Use(gin.Logger()) 
	server.Use(gin.Recovery())


	apiPort := "8001"
	redisAddr := "localhost:6379"
	redisPassword := "my-passwordRedis"
	redisDB := 0


	VotesRepository := repository.NewVotesRepository(redisAddr, redisPassword, redisDB)
	VotesUseCase := usecase.NewVotesUsecase(VotesRepository)
	VotesController := controller.NewVotesController(*VotesUseCase)

	server.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	
	server.GET("/api//health", VotesController.Ping)
	server.POST("/api//votes", VotesController.CreateVotes)
	server.GET("/api/votes/summary", VotesController.GetSummary)


	port := getEnvWithDefault("API_PORT", apiPort)
	log.Printf("Servidor iniciado na porta %s", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}