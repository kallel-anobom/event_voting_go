package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)

type Votes struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Votes     map[string]int `json:"votes" gorm:"type:jsonb"`
	EventName string         `json:"event_name"`
	Date      string         `json:"date"`
	Time      string         `json:"time"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

type VoteStats struct {
	TotalVotes    int64
	VotesByChoice map[string]int64
	VotesByHour   map[string]int64
}

func initDB() {
	var err error
	dsn := "host=db_postgres user=postgres password=postgres123 dbname=db_voting port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	db.AutoMigrate(&Votes{})
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "my-passwordRedis",
		DB:       0,
	})
}

func processVotes() {
	batchSize := 1000
	for {
		votes, err := rdb.LRange(ctx, "vote_queue", 0, int64(batchSize-1)).Result()
		if err != nil || len(votes) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		tx := db.Begin()
		for _, voteData := range votes {
			var vote Votes
			if err := json.Unmarshal([]byte(voteData), &vote); err != nil {
				log.Println("Erro ao desserializar voto: ", err)
				continue
			}

			if err := tx.Create(&vote).Error; err != nil {
				log.Println("Erro ao salvar voto no banco: ", err)
				continue
			}
		}
		tx.Commit()
		rdb.LTrim(ctx, "vote_queue", int64(len(votes)), -1)
	}
}

func main() {
	initDB()
	initRedis()
	log.Println("Worker iniciado... Processando votos...")
	processVotes()
}