package usecase

import (
	"context"
	"fmt"

	"github.com/kallel-anobom/event_voting_go/model"
	"github.com/kallel-anobom/event_voting_go/repository"
)

type VotesUsecase struct {
	repository repository.VotesRepository
}

func NewVotesUsecase(repo repository.VotesRepository) *VotesUsecase {
	return &VotesUsecase{
		repository: repo,
	}
}

func (u *VotesUsecase) CreateVotes(ctx context.Context, votes model.Votes) error {
	exists, err := u.repository.Exists(votes.ID)
	if err != nil {
			return fmt.Errorf("failed to check if votes exist: %w", err)
	}
	if exists {
			return fmt.Errorf("votes already exist")
	}

	return u.repository.CreateVotes(votes)
	}

func (u *VotesUsecase) GetVotes(id int) (*model.Votes, error) {
	intID := id

	votes, err := u.repository.GetVotes(intID)
	if err != nil {
		return nil, err
	}

	return votes, nil
}

func (u *VotesUsecase) GetVotesSummary(ctx context.Context) (*model.VoteSummary, error) {
	return u.repository.GetVotesSummary()
}

func (u *VotesUsecase) Exists(id string) (bool, error) {
	return u.repository.Exists(id)
}

func (u *VotesUsecase) Ping(ctx context.Context) error {
	if err := u.repository.Ping(ctx); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}
	return nil
}