package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kallel-anobom/event_voting_go/api/dto"
)

func (u *VotesUsecase) GetVotesSummary(ctx context.Context) (dto.VoteSummary, error) {
	/*
		Tentar obter o summary do cache Redis
	*/
	var cachedSummary dto.VoteSummary
	err := u.cache.GetJSON(ctx, "votes-summary", &cachedSummary)
	if err == nil && cachedSummary.TotalVotes > 0 {
		fmt.Println("Retornando resumo do cache.")
		return cachedSummary, nil
	}

	/*
		Caso não tenha encontrado no cache, buscar no banco
	*/
	rows, err := u.repository.GetAllVotes()
	if err != nil {
		return dto.VoteSummary{}, err
	}

	/*
		Processar os votos para calcular: Total geral de votos, Total de votos por participante, Total de votos por hora
	*/

	totalVotes := len(rows)
	votesByParticipant := make(map[string]int)
	votesByHour := make(map[string]int)

	loc, _ := time.LoadLocation("America/Sao_Paulo")
	for _, vote := range rows {
		votesByParticipant[strconv.Itoa(vote.ParticipantID)]++

		hour := vote.Date.In(loc).Format("15:00")
		votesByHour[hour]++

	}

	summary := dto.VoteSummary{
		TotalVotes:         totalVotes,
		VotesByParticipant: votesByParticipant,
		VotesByHour:        votesByHour,
	}

	err = u.cache.SetJSON(ctx, "votes-summary", summary, int((5 * time.Minute).Seconds()))
	if err != nil {
		fmt.Println("Erro ao salvar resumo no Redis:", err)
	}

	return summary, nil
}
