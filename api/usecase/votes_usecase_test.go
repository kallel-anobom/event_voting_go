package usecase_test

import (
	"context"
	"testing"
	"time"

	"strconv"

	"github.com/kallel-anobom/event_voting_go/model"
	"github.com/kallel-anobom/event_voting_go/repository"
	"github.com/kallel-anobom/event_voting_go/repository/mocks"
	"github.com/kallel-anobom/event_voting_go/usecase"
	"github.com/stretchr/testify/assert"
)


func TestCreateVotes(t *testing.T) {
	t.Run("Criação bem-sucedida", func(t *testing.T) {
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)
		
		testVotes := model.Votes{
			ID:        "vote-123",
			Name:      "Teste",
			EventName: "Evento Teste",
			Votes:     map[string]int{"op1": 5},
			Date:      time.Now().Format("02/01/2006"),
			Time:      "15:00",
		}

		mockRepo.On("Exists", testVotes.ID).Return(false, nil)
		mockRepo.On("CreateVotes", testVotes).Return(nil)

		err := uc.CreateVotes(context.Background(), testVotes)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro quando votos já existem", func(t *testing.T) {
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)
		
		testVotes := model.Votes{ID: "vote-456"}

		mockRepo.On("Exists", testVotes.ID).Return(true, nil)

		err := uc.CreateVotes(context.Background(), testVotes)
		
		assert.Error(t, err)
		assert.Equal(t, "votes already exist", err.Error())
		mockRepo.AssertExpectations(t)

		mockRepo.AssertNotCalled(t, "CreateVotes")
	})

	t.Run("Erro ao criar votos", func(t *testing.T) {
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)
		
		testVotes := model.Votes{ID: "vote-789"}

		mockRepo.On("Exists", testVotes.ID).Return(false, nil)
		mockRepo.On("CreateVotes", testVotes).Return(assert.AnError)

		err := uc.CreateVotes(context.Background(), testVotes)
		
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetVotes(t *testing.T) {
	t.Run("Obter votos com sucesso", func(t *testing.T) {
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)
		
		expectedVotes := &model.Votes{
			ID:"vote-123", 
			Name:"Teste", 
			Votes:map[string]int{"op1":5}, 
			EventName:"Evento Teste", 
			Date:"25/03/2025", 
			Time:"15:00",
		}

		mockRepo.On("GetVotes", 123).Return(expectedVotes, nil)

		id, _ := strconv.Atoi("123")
		result, err := uc.GetVotes(id)

		assert.NoError(t, err)
		assert.Equal(t, expectedVotes, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro quando votos não existem", func(t *testing.T) {
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)

		mockRepo.On("GetVotes", 456).Return(model.Votes{}, repository.ErrNotFoundVotes)

		id, _ := strconv.Atoi("456")
		_, err := uc.GetVotes(id)

		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFoundVotes, err)
		mockRepo.AssertExpectations(t)
	})
}
func TestGetVotesSummary(t *testing.T) {
	t.Run("Deve retornar resumo de votos com sucesso", func(t *testing.T) {
		// 1. Setup do teste
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)

		// 2. Dados de teste
		expectedSummary := &model.VoteSummary{
			TotalVotes: 100,
			VotesByOption: map[string]int{
				"opcao1": 60,
				"opcao2": 40,
			},
			VotesByHour: map[string]int{
				"10h": 30,
				"11h": 70,
			},
			VotesByParticipant: map[string]int{
				"participante1": 50,
				"participante2": 50,
			},
		}

		// 3. Configuração do mock
		mockRepo.On("GetVotesSummary").Return(expectedSummary, nil)

		// 4. Execução
		summary, err := uc.GetVotesSummary(context.Background())

		// 5. Verificações
		assert.NoError(t, err)
		assert.Equal(t, expectedSummary, summary)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando o repositório falhar", func(t *testing.T) {
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)

		expectedErr := assert.AnError
		mockRepo.On("GetVotesSummary").Return((*model.VoteSummary)(nil), expectedErr)

		summary, err := uc.GetVotesSummary(context.Background())

		assert.Error(t, err)
		assert.Nil(t, summary)
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve lidar com dados vazios corretamente", func(t *testing.T) {
		mockRepo := new(mocks.VotesRepositoryMock)
		uc := usecase.NewVotesUsecase(mockRepo)

		expectedSummary := &model.VoteSummary{
			TotalVotes:        0,
			VotesByOption:     make(map[string]int),
			VotesByHour:       make(map[string]int),
			VotesByParticipant: make(map[string]int),
		}

		mockRepo.On("GetVotesSummary").Return(expectedSummary, nil)

		summary, err := uc.GetVotesSummary(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 0, summary.TotalVotes)
		assert.Empty(t, summary.VotesByOption)
		assert.Empty(t, summary.VotesByHour)
		assert.Empty(t, summary.VotesByParticipant)
		mockRepo.AssertExpectations(t)
	})
}