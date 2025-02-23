package handler

import (
	"net/http"
	"ymb-cloz/internal/service"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	service *service.PlayerService
}

func NewPlayerHandler(service *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func (h *PlayerHandler) GetAllPlayers(c *gin.Context) {
	players, err := h.service.GetAllPlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"players": players})
}

func (h *PlayerHandler) GetTopByWinRate(c *gin.Context) {
	stats, err := h.service.GetTopByWinRate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch win rate statistics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *PlayerHandler) GetTopByGames(c *gin.Context) {
	stats, err := h.service.GetTopByGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch games statistics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *PlayerHandler) GetTopCaptains(c *gin.Context) {
	stats, err := h.service.GetTopCaptains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch captain statistics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *PlayerHandler) GetTopByRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
		return
	}

	stats, err := h.service.GetTopByRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role statistics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
