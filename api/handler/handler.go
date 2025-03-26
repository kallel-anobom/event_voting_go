package handler

import "github.com/kallel-anobom/event_voting_go/api/usecase"

type votesHandler struct {
	votesUsecase usecase.VotesUsecase
}

func NewVotesHandler(usecase usecase.VotesUsecase) *votesHandler {
	return &votesHandler{
		votesUsecase: usecase,
	}
}
