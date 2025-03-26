package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kallel-anobom/event_voting_go/api/model"
)

type VoteRequest struct {
	ParticipantID int `json:"participant_id" binding:"required"`
}

func (v *votesHandler) Vote(ctx *gin.Context) {
	var request VoteRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	v.votesUsecase.Vote(ctx, model.Vote{
		ParticipantID: request.ParticipantID,
		Date:          time.Now(),
	})

	// if len(request.Votes) == 0 {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "O campo 'votes' deve conter pelo menos uma opção",
	// 	})
	// 	return
	// }

	// votes := model.Votes{
	// 	ID:        request.ID,
	// 	Name:      request.Name,
	// 	EventName: request.EventName,
	// 	Votes:     request.Votes,
	// 	Date:      request.Date,
	// 	Time:      request.Time,
	// }

	// if err := v.votesUsecase.Vote(ctx.Request.Context(), votes); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Failed to create votes: " + err.Error(),
	// 	})
	// 	return
	// }

	// exists, err := v.votesUsecase.Exists(request.ID)
	// if err != nil || !exists {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Falha ao verificar a persistência dos dados",
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Vote computed successfully!",
	})
}

// func (v *votesHandler) Ping(ctx *gin.Context) {
// 	start := time.Now()
// 	defer func() {
// 		log.Printf("Health check duration: %v", time.Since(start))
// 	}()

// 	if err := v.votesUsecase.Ping(ctx); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"status": "unhealthy",
// 			"error":  err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"status": "healthy",
// 	})
// }
