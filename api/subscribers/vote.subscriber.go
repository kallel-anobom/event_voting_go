package subscriber

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kallel-anobom/event_voting_go/api/dto"
	"github.com/kallel-anobom/event_voting_go/api/services/pubsub"
)

const VOTE_TOPIC = "votes"

func SubscribeToPubsub(
	pubsub *pubsub.RabbitMQService,
) {
	voteMessage, err := pubsub.GetSubscriber().Subscribe(context.Background(), VOTE_TOPIC)
	if err != nil {
		return
	}

	for {
		select {
		case m := <-voteMessage:
			ProcessVote(m)
		}
	}
}

func ProcessVote(message *message.Message) {
	if message == nil {
		return
	}
	var vote dto.VoteMessage
	if err := json.Unmarshal(message.Payload, &vote); err != nil {
		message.Nack()
		// u.prometheus.AddMetricForErrorInUnmarshalMessage
		return
	}

	// Limpar o cache

	// Se for sucesso
	message.Ack()
	// Se der erro
	// message.Nack()
}
