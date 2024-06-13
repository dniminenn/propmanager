package api

import (
	"net/http"
	"propmanager/internal/app/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats := h.statsService.GetStats()
	c.JSON(http.StatusOK, stats)
}
