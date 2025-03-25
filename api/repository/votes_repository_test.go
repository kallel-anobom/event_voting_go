package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kallel-anobom/event_voting_go/model"
	"github.com/stretchr/testify/assert"
)

func TestVotesRepository(t *testing.T) {
	redisClient := NewRedisClient("localhost:6379", "", 0)
	repo := &votesRepository{redisClient: redisClient}
	ctx := context.Background()

	testVotes := model.Votes{
		ID:        "test1",
		Name:      "Test Event",
		EventName: "Event 1",
		Votes:     map[string]int{"option1": 5},
		Date:      time.Now().Format("02/01/2006"),
		Time:      "15h00",
	}

	redisClient.Client.Del(ctx, "votes:test1")

	t.Run("Create and Get Votes", func(t *testing.T) {
		// Teste Create
		err := repo.CreateVotes(testVotes)
		assert.NoError(t, err)

		// Verificação direta no Redis
		val, err := redisClient.Client.Get(ctx, "votes:test1").Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, val)
	})

	// Teste Exists
	exists, err := repo.Exists(testVotes.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	t.Run("Non-existent Vote", func(t *testing.T) {
		_, err := repo.GetVotes(0)
		assert.Error(t, err)
		assert.Equal(t, "votes not found", err.Error())
	})


	// Limpeza após o teste
	redisClient.Client.Del(ctx, "votes:test1")
}