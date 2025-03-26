package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *votesHandler) GetSummary(ctx *gin.Context) {
	summary, err := c.votesUsecase.GetVotesSummary(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}
