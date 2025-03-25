package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kallel-anobom/event_voting_go/model"
	"github.com/redis/go-redis/v9"
)


type VotesRepository interface {
	CreateVotes(votes model.Votes) error
	GetVotes(id int) (*model.Votes, error)
	GetVotesSummary() (*model.VoteSummary, error)
	Exists(id string) (bool, error)
	Ping(ctx context.Context) error
}

type votesRepository struct {
	redisClient *RedisClient
}

func NewVotesRepository(redisAddr, redisPassword string, redisDB int) VotesRepository {
	redisClient := NewRedisClient(redisAddr, redisPassword, redisDB)
	return &votesRepository{
		redisClient: redisClient,
	}
}


func (vt *votesRepository) CreateVotes(votes model.Votes) error {
	ctx := context.Background()
	
	data, err := json.Marshal(votes)
	if err != nil {
		return errors.New("failed to marshal votes")
	}

	key := "votes:" + votes.ID
	err = vt.redisClient.Client.Set(ctx, key, data, 168*time.Hour).Err()

	if err != nil {
		return errors.New("failed to save votes in Redis")
	}	
	return nil
}

func (vt *votesRepository) GetVotes(id int) (*model.Votes, error) {
	ctx := context.Background()
	key := "votes:" + strconv.Itoa(id)
	data, err := vt.redisClient.Client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFoundVotes
		}
		return nil, errors.New("failed to get votes from Redis")
	}

	var votes model.Votes
	err = json.Unmarshal([]byte(data), &votes)
	if err != nil {
		return nil, errors.New("failed to unmarshal votes")
	}

	return &votes, nil
}

func (r *votesRepository) GetVotesSummary() (*model.VoteSummary, error) {
	ctx := context.Background()
	
	// 1. Obter todas as chaves de votos
	keys, err := r.redisClient.Client.Keys(ctx, "votes:*").Result()
	if err != nil {
			return nil, err
	}

	summary := &model.VoteSummary{
			VotesByOption:     make(map[string]int),
			VotesByHour:       make(map[string]int),
			VotesByParticipant: make(map[string]int),
	}

	// 2. Processar cada voto
	for _, key := range keys {
			data, err := r.redisClient.Client.Get(ctx, key).Result()
			if err != nil {
					continue // Ou tratar o erro conforme sua política
			}

			var vote model.Votes
			if err := json.Unmarshal([]byte(data), &vote); err != nil {
					continue
			}

			// 3. Consolidar dados
			summary.TotalVotes += len(vote.Votes)
			
			// Por opção
			for option, count := range vote.Votes {
					summary.VotesByOption[option] += count
			}
			
			// Por hora (assumindo que vote.Time está no formato "15:00")
			hour := strings.Split(vote.Time, ":")[0] + "h"
			summary.VotesByHour[hour] += len(vote.Votes)
			
			// Por participante
			summary.VotesByParticipant[vote.Name] += len(vote.Votes)
	}

	return summary, nil
}

func (r *votesRepository) Exists(id string) (bool, error) {
	ctx := context.Background()
	key := "votes:" + id
	
	exists, err := r.redisClient.Client.Exists(ctx, key).Result()
	if err != nil {
			return false, err
	}
	
	return exists == 1, nil
}
func (r *votesRepository) Ping(ctx context.Context) error {
	_, err := r.redisClient.Client.Ping(ctx).Result()
	return err
}