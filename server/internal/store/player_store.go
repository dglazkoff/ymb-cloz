package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type PlayerStore struct {
	db *sql.DB
}

func NewPlayerStore(db *sql.DB) *PlayerStore {
	return &PlayerStore{db: db}
}

func (s *PlayerStore) GetAllPlayers() ([]Player, error) {
	query := `SELECT id, nickname, COALESCE(games_played, ARRAY[]::UUID[]) FROM players`
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("error querying players: %v", err)
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var player Player
		var gamesPlayed []sql.NullString
		if err := rows.Scan(&player.ID, &player.Nickname, pq.Array(&gamesPlayed)); err != nil {
			log.Printf("error scanning player: %v", err)
			return nil, err
		}

		player.GamesPlayed = make([]string, 0, len(gamesPlayed))
		for _, g := range gamesPlayed {
			if g.Valid {
				player.GamesPlayed = append(player.GamesPlayed, g.String)
			}
		}

		players = append(players, player)
	}
	if err = rows.Err(); err != nil {
		log.Printf("error iterating players: %v", err)
		return nil, err
	}

	return players, nil
}

type Player struct {
	ID          string   `json:"id"`
	Nickname    string   `json:"nickname"`
	GamesPlayed []string `json:"games_played"`
}

type PlayerStats struct {
	ID       string
	Nickname string
	Stats    string
}

func (s *PlayerStore) GetTopByWinRate() ([]PlayerStats, error) {
	query := `
		SELECT 
			p.id,
			p.nickname,
			CAST(COUNT(CASE WHEN g.is_winner = true THEN 1 END) AS float) / CAST(COUNT(*) AS float) * 100 as winrate,
			COUNT(CASE WHEN g.is_winner = true THEN 1 END) as wins,
			COUNT(*) as total_games
		FROM players p
		JOIN game_players g ON p.id = g.player_id
		GROUP BY p.id, p.nickname
		HAVING COUNT(*) > 0
		ORDER BY winrate DESC
		LIMIT 10`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PlayerStats
	for rows.Next() {
		var stat PlayerStats
		var winrate float64
		var wins, totalGames int
		if err := rows.Scan(&stat.ID, &stat.Nickname, &winrate, &wins, &totalGames); err != nil {
			return nil, err
		}
		stat.Stats = fmt.Sprintf("%.1f%% (%d/%d)", winrate, wins, totalGames)
		stats = append(stats, stat)
	}
	return stats, rows.Err()
}

func (s *PlayerStore) GetTopByGames() ([]PlayerStats, error) {
	query := `
		SELECT 
			p.id,
			p.nickname,
			COUNT(*) as games
		FROM players p
		JOIN game_players g ON p.id = g.player_id
		GROUP BY p.id, p.nickname
		ORDER BY games DESC
		LIMIT 10`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PlayerStats
	for rows.Next() {
		var stat PlayerStats
		var games int
		if err := rows.Scan(&stat.ID, &stat.Nickname, &games); err != nil {
			return nil, err
		}
		stat.Stats = fmt.Sprintf("%d games", games)
		stats = append(stats, stat)
	}
	return stats, rows.Err()
}

func (s *PlayerStore) GetTopCaptains() ([]PlayerStats, error) {
	query := `
		SELECT 
			p.id,
			p.nickname,
			CAST(COUNT(CASE WHEN g.is_winner = true THEN 1 END) AS float) / CAST(COUNT(*) AS float) * 100 as winrate,
			COUNT(CASE WHEN g.is_winner = true THEN 1 END) as wins,
			COUNT(*) as total_games
		FROM players p
		JOIN game_players g ON p.id = g.player_id
		WHERE g.is_captain = true
		GROUP BY p.id, p.nickname
		HAVING COUNT(*) > 0
		ORDER BY winrate DESC
		LIMIT 10`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PlayerStats
	for rows.Next() {
		var stat PlayerStats
		var winrate float64
		var wins, totalGames int
		if err := rows.Scan(&stat.ID, &stat.Nickname, &winrate, &wins, &totalGames); err != nil {
			return nil, err
		}
		stat.Stats = fmt.Sprintf("%.1f%% (%d/%d)", winrate, wins, totalGames)
		stats = append(stats, stat)
	}
	return stats, rows.Err()
}

func (s *PlayerStore) GetTopByRole(role string) ([]PlayerStats, error) {
	query := `
		SELECT 
			p.id,
			p.nickname,
			CAST(COUNT(CASE WHEN g.is_winner = true THEN 1 END) AS float) / CAST(COUNT(*) AS float) * 100 as winrate,
			COUNT(CASE WHEN g.is_winner = true THEN 1 END) as wins,
			COUNT(*) as total_games
		FROM players p
		JOIN game_players g ON p.id = g.player_id
		WHERE g.role = $1
		GROUP BY p.id, p.nickname
		HAVING COUNT(*) > 0
		ORDER BY winrate DESC
		LIMIT 10`

	rows, err := s.db.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PlayerStats
	for rows.Next() {
		var stat PlayerStats
		var winrate float64
		var wins, totalGames int
		if err := rows.Scan(&stat.ID, &stat.Nickname, &winrate, &wins, &totalGames); err != nil {
			return nil, err
		}
		stat.Stats = fmt.Sprintf("%.1f%% (%d/%d)", winrate, wins, totalGames)
		stats = append(stats, stat)
	}
	return stats, rows.Err()
}
