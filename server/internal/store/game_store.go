package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type GameStore interface {
	BeginTx() (*sql.Tx, error)
	CreateGameTx(tx *sql.Tx, game *Game) error
	GetOrCreatePlayerTx(tx *sql.Tx, nickname string) (string, error)
	GetPlayerByIDTx(tx *sql.Tx, id string) (bool, error)
	CreateGamePlayersTx(tx *sql.Tx, gameID string, players []GamePlayer) error
	UpdatePlayersGamesTx(tx *sql.Tx, gameID string, playerIDs []string) error
}

type PostgresGameStore struct {
	db *sql.DB
}

func NewGameStore(db *sql.DB) GameStore {
	return &PostgresGameStore{db: db}
}

type Game struct {
	ID        string
	Timestamp string
	Winner    string
}

type GamePlayer struct {
	GameID    string
	PlayerID  string
	Team      string
	Role      string
	IsCaptain bool
	IsWinner  bool
}

func (s *PostgresGameStore) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *PostgresGameStore) GetOrCreatePlayerTx(tx *sql.Tx, nickname string) (string, error) {
	var playerID string

	// Try to find existing player
	err := tx.QueryRow("SELECT id FROM players WHERE nickname = $1", nickname).Scan(&playerID)
	if err == nil {
		// Player found
		return playerID, nil
	}

	if err != sql.ErrNoRows {
		log.Printf("error checking player existence: %v", err)
		return "", fmt.Errorf("error checking player existence: %v", err)
	}

	// Player not found, create new one
	err = tx.QueryRow(`
		INSERT INTO players (nickname)
		VALUES ($1)
		RETURNING id`, nickname).Scan(&playerID)
	if err != nil {
		log.Printf("error creating player: %v", err)
		return "", fmt.Errorf("error creating player: %v", err)
	}

	return playerID, nil
}

func (s *PostgresGameStore) GetPlayerByIDTx(tx *sql.Tx, id string) (bool, error) {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM players WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking player existence by ID: %v", err)
	}
	return exists, nil
}

func (s *PostgresGameStore) CreateGameTx(tx *sql.Tx, game *Game) error {
	query := `
		INSERT INTO games (winner)
		VALUES ($1)
		RETURNING id, timestamp`

	err := tx.QueryRow(query, game.Winner).Scan(&game.ID, &game.Timestamp)
	if err != nil {
		return fmt.Errorf("error creating game: %v", err)
	}

	return nil
}

func (s *PostgresGameStore) CreateGamePlayersTx(tx *sql.Tx, gameID string, players []GamePlayer) error {
	query := `
		INSERT INTO game_players (game_id, player_id, team, role, is_captain, is_winner)
		VALUES ($1, $2, $3, $4, $5, $6)`

	for _, player := range players {
		_, err := tx.Exec(query, gameID, player.PlayerID, player.Team, player.Role, player.IsCaptain, player.IsWinner)
		if err != nil {
			return fmt.Errorf("error creating game player: %v", err)
		}
	}

	return nil
}

func (s *PostgresGameStore) UpdatePlayersGamesTx(tx *sql.Tx, gameID string, playerIDs []string) error {
	query := `
		UPDATE players 
		SET games_played = array_append(games_played, $1)
		WHERE id = ANY($2)`

	_, err := tx.Exec(query, gameID, pq.Array(playerIDs))
	if err != nil {
		return fmt.Errorf("error updating players games count: %v", err)
	}

	return nil
}
