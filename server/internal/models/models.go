package models

import (
	"time"
)

type Player struct {
	ID       string   `json:"id"`
	Nickname string   `json:"nickname"`
	Games    []string `json:"games"` // Array of game IDs
}

type GamePlayer struct {
	GameID    string `json:"game_id"`
	PlayerID  string `json:"player_id"`
	Team      string `json:"team"`
	Role      Role   `json:"role"`
	IsCaptain bool   `json:"is_captain"`
	IsWinner  bool   `json:"is_winner"`
}

type Role string

const (
	Carry   Role = "carry"
	Mid     Role = "mid"
	Offlane Role = "offlane"
	Pos4    Role = "pos4"
	Pos5    Role = "pos5"
)

type Game struct {
	ID          string       `json:"id"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
	RadiantTeam []GamePlayer `json:"radiant_team"`
	DireTeam    []GamePlayer `json:"dire_team"`
	Winner      string       `json:"winner"`
}
