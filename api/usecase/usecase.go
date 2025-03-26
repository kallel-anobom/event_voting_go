package usecase

import (
	"github.com/kallel-anobom/event_voting_go/api/repository"
	"github.com/kallel-anobom/event_voting_go/api/services/cache"
	"github.com/kallel-anobom/event_voting_go/api/services/pubsub"
)

type VotesUsecase struct {
	repository repository.VotesRepository
	pubsub     *pubsub.RabbitMQService
	cache      *cache.RedisService
}

func NewVotesUsecase(repo repository.VotesRepository, redis *cache.RedisService, pubsub *pubsub.RabbitMQService) VotesUsecase {
	return VotesUsecase{
		repository: repo,
		pubsub:     pubsub,
		cache:      redis,
	}
}
