package subscriber

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kallel-anobom/event_voting_go/api/dto"
	"github.com/kallel-anobom/event_voting_go/api/model"
	"github.com/kallel-anobom/event_voting_go/api/repository"
	"github.com/kallel-anobom/event_voting_go/api/services/cache"
	"github.com/kallel-anobom/event_voting_go/api/services/metrics"
	"github.com/kallel-anobom/event_voting_go/api/services/pubsub"
	"github.com/prometheus/client_golang/prometheus"
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
	metrics.MessagesReceived.Inc()

	timer := prometheus.NewTimer(metrics.MessageProcessingTime)
	defer timer.ObserveDuration()

	if message == nil {
		return
	}
	var vote dto.VoteMessage
	if err := json.Unmarshal(message.Payload, &vote); err != nil {
		message.Nack()
		metrics.MessagesFailed.WithLabelValues("unmarshal_error").Inc()
		return
	}

	err := vs.repository.AddVote(model.Vote{
		ParticipantID: vote.ParticipantID,
	})
	if err != nil {
		message.Nack()
		metrics.MessagesFailed.WithLabelValues("db_error").Inc()
		return
	}

	vs.cache.Client.Del(context.Background(), "votes-summary")
	metrics.CacheInvalidations.Inc()

	metrics.MessagesProcessedSuccessfully.Inc()
	message.Ack()
}
