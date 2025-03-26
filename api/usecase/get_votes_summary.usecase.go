package usecase

import (
	"context"
	"encoding/json"
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

	for _, vote := range rows {

		votesByParticipant[strconv.Itoa(vote.ParticipantID)]++

		hour := vote.Date.Format("15:00")
		votesByHour[hour]++
	}

	summary := dto.VoteSummary{
		TotalVotes:         totalVotes,
		VotesByParticipant: votesByParticipant,
		VotesByHour:        votesByHour,
	}

	data, _ := json.Marshal(summary)
	err = u.cache.SetJSON(ctx, "votes-summary", string(data), int((5 * time.Minute).Seconds()))
	if err != nil {
		fmt.Println("Erro ao salvar resumo no Redis:", err)
	}

	return summary, nil
	// Acessar o REDIS pra ver se tem já o summary cacheado
	// Se não tiver, busca TODOS OS VOTOS DO BANCO (tem que ver se o mongo tem algum método de agregar e já trazer a média do numero de votos por participante)
	// Processa a média dos votos AQUI
	// E salve no cache AQUI

	// u.cache.Get("votes-summary")
	// Se existir no cache, retorna do cache

	// Se não existir no cache, busca no banco
	// rows, err := u.repository.GetAllVotes()
	// if err != nil {
	// 	return dto.VoteSummary{}, err
	// }

	// fmt.Println(rows)

	// AQUI VOCÊ PROCESSA AS ROWS PARA GERAR O OBJETO `model`

	// ctx := context.Background()

	// // 1. Obter todas as chaves de votos
	// keys, err := r.mongoService.Client.Keys(ctx, "votes:*").Result()
	// if err != nil {
	// 	return nil, err
	// }

	// summary := &model.VoteSummary{
	// 	VotesByOption:      make(map[string]int),
	// 	VotesByHour:        make(map[string]int),
	// 	VotesByParticipant: make(map[string]int),
	// }

	// // 2. Processar cada voto
	// for _, key := range keys {
	// 	data, err := r.mongoService.Client.Get(ctx, key).Result()
	// 	if err != nil {
	// 		continue // Ou tratar o erro conforme sua política
	// 	}

	// 	var vote model.Votes
	// 	if err := json.Unmarshal([]byte(data), &vote); err != nil {
	// 		continue
	// 	}

	// 	// 3. Consolidar dados
	// 	summary.TotalVotes += len(vote.Votes)

	// 	// Por opção
	// 	for option, count := range vote.Votes {
	// 		summary.VotesByOption[option] += count
	// 	}

	// 	// Por hora (assumindo que vote.Time está no formato "15:00")
	// 	hour := strings.Split(vote.Time, ":")[0] + "h"
	// 	summary.VotesByHour[hour] += len(vote.Votes)

	// 	// Por participante
	// 	summary.VotesByParticipant[vote.Name] += len(vote.Votes)
	// }

	// AO FINAL ANTES DE RETORNAR O OBJETO DE SUMMARY, SALVA NO REDIS
	// PARA QUE NA PRÓXIMA VEZ QUE ALGUÉM REQUISITAR O SUMMARY, ELE
	// JÁ ESTEJA CACHEADO

	// return summary, nil
}
