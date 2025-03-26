package usecase

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kallel-anobom/event_voting_go/api/model"
	"github.com/kallel-anobom/event_voting_go/api/subscriber"
)

func (u *VotesUsecase) Vote(ctx context.Context, vote model.Vote) error {
	jsonData, err := json.Marshal(vote)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), jsonData)

	if err := u.pubsub.GetPublisher().Publish(subscriber.VOTE_TOPIC, msg); err != nil {
		return err
	}

	return nil
}
