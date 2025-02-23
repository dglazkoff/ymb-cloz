export type CreateGameRequest = {
	radiant_players: GamePlayerInput[];
	dire_players:    GamePlayerInput[];
	winner: string;
}

export type GamePlayerInput = {
	id?:        string;
	nickname?:  string;
	role:      string;
	is_captain: boolean;
}