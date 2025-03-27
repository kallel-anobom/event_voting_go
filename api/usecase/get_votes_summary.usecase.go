package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kallel-anobom/event_voting_go/api/dto"
	"github.com/kallel-anobom/event_voting_go/api/services/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func (u *VotesUsecase) GetVotesSummary(ctx context.Context) (dto.VoteSummary, error) {
	timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues("GetVotesSummary"))
	defer timer.ObserveDuration()

	/*
		Tentar obter o summary do cache Redis
	*/
	var cachedSummary dto.VoteSummary
	err := u.cache.GetJSON(ctx, "votes-summary", &cachedSummary)
	if err == nil && cachedSummary.TotalVotes > 0 {
		fmt.Println("Retornando resumo do cache.")
		metrics.TotalVotes.Add(float64(cachedSummary.TotalVotes))
		return cachedSummary, nil
	}

	/*
		Caso não tenha encontrado no cache, buscar no banco
	*/
	rows, err := u.repository.GetAllVotes()
	if err != nil {
		metrics.ErrorsCount.WithLabelValues("database_error").Inc()
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
		participantID := strconv.Itoa(vote.ParticipantID)

		votesByParticipant[participantID]++
		metrics.VotesByParticipant.WithLabelValues(participantID).Inc()

		hour := vote.Date.In(loc).Format("15:00")
		votesByHour[hour]++
	}

	metrics.TotalVotes.Add(float64(totalVotes))

	summary := dto.VoteSummary{
		TotalVotes:         totalVotes,
		VotesByParticipant: votesByParticipant,
		VotesByHour:        votesByHour,
	}

	err = u.cache.SetJSON(ctx, "votes-summary", summary, int((5 * time.Minute).Seconds()))
	if err != nil {
		metrics.ErrorsCount.WithLabelValues("cache_error").Inc()
		fmt.Println("Erro ao salvar resumo no Redis:", err)
	}

	return summary, nil
}
