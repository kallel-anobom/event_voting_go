package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kallel-anobom/event_voting_go/model"
	"github.com/kallel-anobom/event_voting_go/usecase"
)

type votesController struct{
	VotesUseCase usecase.VotesUsecase
}

func NewVotesController(usecase usecase.VotesUsecase) *votesController {
	return &votesController {
		VotesUseCase: usecase,
	}
}

type CreateVotesRequest struct {
	ID      	string         `json:"id" binding:"required"`
	Name			string `json:"name" binding:"required"`
	EventName string         `json:"event_name" binding:"required"`
	Votes   map[string]int `json:"votes" binding:"required"`
	Date	string `json:"date"`
	Time	string `json:"time"`
}

func (v *votesController) CreateVotes(ctx *gin.Context) {
	var request CreateVotesRequest
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	if len(request.Votes) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "O campo 'votes' deve conter pelo menos uma opção",
		})
		return
}

	votes := model.Votes{
		ID: request.ID,
		Name: request.Name,
		EventName: request.EventName,
		Votes: request.Votes,
		Date: request.Date,
		Time: request.Time,
	}

	if err := v.VotesUseCase.CreateVotes(ctx.Request.Context(), votes); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create votes: " + err.Error(),
		})
		return
	}


	exists, err := v.VotesUseCase.Exists(request.ID)
	if err != nil || !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Falha ao verificar a persistência dos dados",
			})
			return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Votes created successfully!",
		"data": votes,
	})
}

func (c *votesController) GetSummary(ctx *gin.Context) {
	summary, err := c.VotesUseCase.GetVotesSummary(ctx.Request.Context())
	if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}
	
	ctx.JSON(http.StatusOK, summary)
}

func (v *votesController) Ping(ctx *gin.Context) {
	start := time.Now()
	defer func() {
			log.Printf("Health check duration: %v", time.Since(start))
	}()
	
	if err := v.VotesUseCase.Ping(ctx); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": "unhealthy",
					"error": err.Error(),
			})
			return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
			"status": "healthy",
	})
}