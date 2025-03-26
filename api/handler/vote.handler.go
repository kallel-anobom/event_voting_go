package handler

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kallel-anobom/event_voting_go/api/model"
)

type VoteRequest struct {
	ParticipantID int `json:"participant_id" binding:"required"`
}

func (v *votesHandler) Vote(ctx *gin.Context) {
	userAgent := ctx.GetHeader("User-Agent")

	var request VoteRequest

	if !isHumanVote(userAgent) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Apenas humanos podem votar"})
		return
	}

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

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Vote computed successfully!",
	})
}

func isHumanVote(userAgent string) bool {
	botRegex := regexp.MustCompile(`(?i)bot|crawler|spider|curl|wget`)
	return !botRegex.MatchString(userAgent)
}
