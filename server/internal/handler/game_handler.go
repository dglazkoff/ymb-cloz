package handler

import (
	"net/http"
	"ymb-cloz/internal/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	service service.GameService
}

func NewGameHandler(service service.GameService) *GameHandler {
	return &GameHandler{service: service}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var req service.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if len(req.RadiantPlayers) != 5 || len(req.DirePlayers) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "each team must have exactly 5 players"})
		return
	}

	if req.Winner != "RADIANT" && req.Winner != "DIRE" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "winner must be either RADIANT or DIRE"})
		return
	}

	// Validate roles and captains
	radiantCaptains := 0
	direCaptains := 0
	roles := make(map[string]bool)

	// Check Radiant team
	for _, p := range req.RadiantPlayers {
		if p.ID != nil && p.Nickname != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "player ID and nickname cannot both be provided"})
			return
		}
		if p.ID == nil && p.Nickname == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "player ID or nickname must be provided"})
			return
		}
		if p.Role != "carry" && p.Role != "mid" && p.Role != "offlane" && p.Role != "pos4" && p.Role != "pos5" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role: " + p.Role})
			return
		}
		if p.IsCaptain {
			radiantCaptains++
		}
		roles[p.Role] = true
	}
	if radiantCaptains != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Radiant team must have exactly one captain"})
		return
	}
	if len(roles) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Radiant team must have all unique roles"})
		return
	}

	// Reset roles map for Dire team
	roles = make(map[string]bool)

	// Check Dire team
	for _, p := range req.DirePlayers {
		if p.ID != nil && p.Nickname != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "player ID and nickname cannot both be provided"})
			return
		}
		if p.ID == nil && p.Nickname == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "player ID or nickname must be provided"})
			return
		}
		if p.Role != "carry" && p.Role != "mid" && p.Role != "offlane" && p.Role != "pos4" && p.Role != "pos5" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role: " + p.Role})
			return
		}
		if p.IsCaptain {
			direCaptains++
		}
		roles[p.Role] = true
	}
	if direCaptains != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dire team must have exactly one captain"})
		return
	}
	if len(roles) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dire team must have all unique roles"})
		return
	}

	err := h.service.CreateGame(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "game created successfully"})
}
