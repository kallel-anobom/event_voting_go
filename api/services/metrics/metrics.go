package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalVotes = promauto.NewCounter(prometheus.CounterOpts{
		Name: "votes_total",
		Help: "Número total de votos registrados",
	})

	VotesByParticipant = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "votes_by_participant_total",
			Help: "Número de votos por participante",
		},
		[]string{"participant_id"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"endpoint"},
	)

	ErrorsCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total de erros ocorridos",
		},
		[]string{"type"},
	)

	MessagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pubsub_messages_received_total",
		Help: "Total de mensagens recebidas do pub/sub",
	})

	MessagesProcessedSuccessfully = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pubsub_messages_processed_successfully_total",
		Help: "Total de mensagens processadas com sucesso",
	})

	MessagesFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pubsub_messages_failed_total",
			Help: "Total de mensagens que falharam no processamento",
		},
		[]string{"error_type"},
	)

	MessageProcessingTime = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "pubsub_message_processing_time_seconds",
			Help:    "Tempo para processar uma mensagem",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2}, // Tempos em segundos
		},
	)

	CacheInvalidations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_invalidations_total",
		Help: "Total de vezes que o cache foi limpo",
	})
)
