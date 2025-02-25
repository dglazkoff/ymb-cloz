package service

import (
	"ymb-cloz/internal/store"
)

type PlayerService struct {
	store *store.PlayerStore
}

func NewPlayerService(store *store.PlayerStore) *PlayerService {
	return &PlayerService{store: store}
}

func (s *PlayerService) GetAllPlayers() ([]store.Player, error) {
	return s.store.GetAllPlayers()
}

func (s *PlayerService) GetTopByWinRate() ([]store.PlayerStats, error) {
	return s.store.GetTopByWinRate()
}

func (s *PlayerService) GetTopByGames() ([]store.PlayerStats, error) {
	return s.store.GetTopByGames()
}

func (s *PlayerService) GetTopCaptains() ([]store.PlayerStats, error) {
	return s.store.GetTopCaptains()
}

func (s *PlayerService) GetTopByRole(role string) ([]store.PlayerStats, error) {
	return s.store.GetTopByRole(role)
}

// test
func (s *PlayerService) GetProkurorStats() (store.PlayerStats, error) {
	return s.store.GetPlayerStats("9cbeb686-ff5f-4c58-bd66-1c0abd54f187")
}
