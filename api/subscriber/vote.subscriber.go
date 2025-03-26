package subscriber

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kallel-anobom/event_voting_go/api/dto"
	"github.com/kallel-anobom/event_voting_go/api/model"
	"github.com/kallel-anobom/event_voting_go/api/repository"
	"github.com/kallel-anobom/event_voting_go/api/services/cache"
	"github.com/kallel-anobom/event_voting_go/api/services/pubsub"
)

const VOTE_TOPIC = "votes"

type voteSubscriber struct {
	repository repository.VotesRepository
	cache      *cache.RedisService
	pubsub     *pubsub.RabbitMQService
}

func NewVoteSubscribers(
	repository repository.VotesRepository,
	cache *cache.RedisService,
	pubsub *pubsub.RabbitMQService,
) voteSubscriber {
	return voteSubscriber{repository, cache, pubsub}
}

func (vs voteSubscriber) SubscribeToPubsub() {
	voteMessage, err := vs.pubsub.GetSubscriber().Subscribe(context.Background(), VOTE_TOPIC)
	if err != nil {
		return
	}

	for m := range voteMessage {
		vs.processVote(m)
	}
}

func (vs voteSubscriber) processVote(message *message.Message) {
	if message == nil {
		return
	}
	var vote dto.VoteMessage
	if err := json.Unmarshal(message.Payload, &vote); err != nil {
		message.Nack()
		// u.prometheus.AddMetricForErrorInUnmarshalMessage
		return
	}

	err := vs.repository.AddVote(model.Vote{
		ParticipantID: vote.ParticipantID,
	})
	if err != nil {
		message.Nack()
		// u.prometheus.AddMetricForErrorInAddVote
		return
	}

	// Limpar o cache do votes-summary
	vs.cache.Client.Del(context.Background(), "votes-summary")

	// Ack a mensagem
	message.Ack()
}
