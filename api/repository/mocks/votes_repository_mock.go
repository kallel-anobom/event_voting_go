// repository/mocks/votes_repository.go
package mocks

import (
	"context"

	"github.com/kallel-anobom/event_voting_go/model"
	"github.com/stretchr/testify/mock"
)

type VotesRepositoryMock struct {
	mock.Mock
}

func (m *VotesRepositoryMock) CreateVotes(votes model.Votes) error {
	args := m.Called(votes)
	return args.Error(0)
}

func (m *VotesRepositoryMock) GetVotes(id int) (*model.Votes, error) {
	return m.Called(id).Get(0).(*model.Votes), m.Called(id).Error(1)
}

func (m *VotesRepositoryMock) GetVotesSummary() (*model.VoteSummary, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VoteSummary), args.Error(1)
}


func (m *VotesRepositoryMock) Exists(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *VotesRepositoryMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
