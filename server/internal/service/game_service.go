package service

import (
	"database/sql"
	"fmt"

	"ymb-cloz/internal/store"
)

type GameService interface {
	CreateGame(req *CreateGameRequest) error
}

type gameService struct {
	store store.GameStore
}

func NewGameService(store store.GameStore) GameService {
	return &gameService{store: store}
}

type CreateGameRequest struct {
	RadiantPlayers []GamePlayerInput `json:"radiant_players"`
	DirePlayers    []GamePlayerInput `json:"dire_players"`
	Winner         string            `json:"winner"`
}

type GamePlayerInput struct {
	Nickname  *string `json:"nickname"`
	ID        *string `json:"id"`
	Role      string  `json:"role"`
	IsCaptain bool    `json:"is_captain"`
}

func (s *gameService) getPlayerID(tx *sql.Tx, input GamePlayerInput) (string, error) {
	// If ID is provided, verify it exists
	if input.ID != nil {
		exists, err := s.store.GetPlayerByIDTx(tx, *input.ID)
		if err != nil {
			return "", fmt.Errorf("error checking player ID: %v", err)
		}
		if !exists {
			return "", fmt.Errorf("player with ID %s not found", *input.ID)
		}
		return *input.ID, nil
	}

	// If nickname is provided, get or create player
	if input.Nickname != nil {
		playerID, err := s.store.GetOrCreatePlayerTx(tx, *input.Nickname)
		if err != nil {
			return "", fmt.Errorf("error getting/creating player by nickname: %v", err)
		}
		return playerID, nil
	}

	return "", fmt.Errorf("either player ID or nickname must be provided")
}

func (s *gameService) CreateGame(req *CreateGameRequest) error {
	// Create game record
	game := &store.Game{
		Winner: req.Winner,
	}

	// Begin transaction
	tx, err := s.store.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Create game
	err = s.store.CreateGameTx(tx, game)
	if err != nil {
		return fmt.Errorf("failed to create game: %v", err)
	}

	// Prepare players data
	var players []store.GamePlayer

	// Add Radiant players
	for _, p := range req.RadiantPlayers {
		playerID, err := s.getPlayerID(tx, p)
		if err != nil {
			return fmt.Errorf("failed to process Radiant player: %v", err)
		}

		players = append(players, store.GamePlayer{
			GameID:    game.ID,
			PlayerID:  playerID,
			Team:      "RADIANT",
			Role:      p.Role,
			IsCaptain: p.IsCaptain,
			IsWinner:  game.Winner == "RADIANT",
		})
	}

	// Add Dire players
	for _, p := range req.DirePlayers {
		playerID, err := s.getPlayerID(tx, p)
		if err != nil {
			return fmt.Errorf("failed to process Dire player: %v", err)
		}

		players = append(players, store.GamePlayer{
			GameID:    game.ID,
			PlayerID:  playerID,
			Team:      "DIRE",
			Role:      p.Role,
			IsCaptain: p.IsCaptain,
			IsWinner:  game.Winner == "DIRE",
		})
	}

	// Create game players
	err = s.store.CreateGamePlayersTx(tx, game.ID, players)
	if err != nil {
		return fmt.Errorf("failed to create game players: %v", err)
	}

	// Update games_played for all players
	playerIDs := make([]string, len(players))
	for i, player := range players {
		playerIDs[i] = player.PlayerID
	}

	err = s.store.UpdatePlayersGamesTx(tx, game.ID, playerIDs)
	if err != nil {
		return fmt.Errorf("failed to update players games count: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
